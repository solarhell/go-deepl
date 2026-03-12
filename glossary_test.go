package deepl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/solarhell/go-deepl"
)

func glossaryHandler(t *testing.T, method, pathSuffix string, statusCode int, body string, onRequest func(*http.Request)) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Fatalf("expected %s, got %s", method, r.Method)
		}
		wantPath := "/glossaries" + pathSuffix
		if r.URL.Path != wantPath {
			t.Fatalf("expected %s, got %s", wantPath, r.URL.Path)
		}
		if onRequest != nil {
			onRequest(r)
		}
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

func TestCreateGlossary(t *testing.T) {
	var captured *http.Request
	server := newTestServer(t, glossaryHandler(t, "POST", "", http.StatusCreated,
		`{"glossary_id":"gl-123","name":"My Glossary","ready":true,"source_lang":"EN","target_lang":"DE","entry_count":2}`,
		func(r *http.Request) {
			r.ParseForm()
			captured = r
		},
	))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	glossary, err := client.CreateGlossary(context.Background(), "My Glossary",
		deepl.English, deepl.German,
		[]deepl.GlossaryEntry{
			{Source: "hello", Target: "hallo"},
			{Source: "world", Target: "Welt"},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("request was not captured")
	}

	if got := captured.FormValue("name"); got != "My Glossary" {
		t.Errorf("name = %q, want %q", got, "My Glossary")
	}
	if got := captured.FormValue("source_lang"); got != "EN" {
		t.Errorf("source_lang = %q, want %q", got, "EN")
	}
	if got := captured.FormValue("target_lang"); got != "DE" {
		t.Errorf("target_lang = %q, want %q", got, "DE")
	}
	if got := captured.FormValue("entries_format"); got != "tsv" {
		t.Errorf("entries_format = %q, want %q", got, "tsv")
	}
	if got := captured.FormValue("entries"); got != "hello\thallo\nworld\tWelt" {
		t.Errorf("entries = %q, want %q", got, "hello\thallo\nworld\tWelt")
	}

	if glossary.GlossaryID != "gl-123" {
		t.Errorf("GlossaryID = %q, want %q", glossary.GlossaryID, "gl-123")
	}
	if glossary.Name != "My Glossary" {
		t.Errorf("Name = %q, want %q", glossary.Name, "My Glossary")
	}
	if !glossary.Ready {
		t.Error("Ready = false, want true")
	}
	if glossary.EntryCount != 2 {
		t.Errorf("EntryCount = %d, want %d", glossary.EntryCount, 2)
	}
}

func TestCreateGlossary_error(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "POST", "", http.StatusBadRequest, `{"message":"bad request"}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	_, err := client.CreateGlossary(context.Background(), "test",
		deepl.English, deepl.German,
		[]deepl.GlossaryEntry{{Source: "a", Target: "b"}},
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var deeplError deepl.Error
	if !errors.As(err, &deeplError) {
		t.Fatalf("expected deepl.Error, got %T", err)
	}
	if deeplError.Code != http.StatusBadRequest {
		t.Errorf("Code = %d, want %d", deeplError.Code, http.StatusBadRequest)
	}
}

func TestListGlossaries(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "GET", "", http.StatusOK,
		`{"glossaries":[{"glossary_id":"gl-1","name":"G1","ready":true,"source_lang":"EN","target_lang":"DE","entry_count":1},{"glossary_id":"gl-2","name":"G2","ready":false,"source_lang":"FR","target_lang":"DE","entry_count":3}]}`,
		nil,
	))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	glossaries, err := client.ListGlossaries(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(glossaries) != 2 {
		t.Fatalf("len(glossaries) = %d, want 2", len(glossaries))
	}
	if glossaries[0].GlossaryID != "gl-1" {
		t.Errorf("glossaries[0].GlossaryID = %q, want %q", glossaries[0].GlossaryID, "gl-1")
	}
	if glossaries[1].Name != "G2" {
		t.Errorf("glossaries[1].Name = %q, want %q", glossaries[1].Name, "G2")
	}
}

func TestListGlossaries_empty(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "GET", "", http.StatusOK, `{"glossaries":[]}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	glossaries, err := client.ListGlossaries(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(glossaries) != 0 {
		t.Errorf("len(glossaries) = %d, want 0", len(glossaries))
	}
}

func TestListGlossary(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "GET", "/gl-123", http.StatusOK,
		`{"glossary_id":"gl-123","name":"Test","ready":true,"source_lang":"EN","target_lang":"DE","entry_count":5}`,
		nil,
	))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	glossary, err := client.ListGlossary(context.Background(), "gl-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if glossary.GlossaryID != "gl-123" {
		t.Errorf("GlossaryID = %q, want %q", glossary.GlossaryID, "gl-123")
	}
	if glossary.EntryCount != 5 {
		t.Errorf("EntryCount = %d, want %d", glossary.EntryCount, 5)
	}
}

func TestListGlossary_notFound(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "GET", "/nonexistent", http.StatusNotFound, `{"message":"not found"}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	_, err := client.ListGlossary(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var deeplError deepl.Error
	if !errors.As(err, &deeplError) {
		t.Fatalf("expected deepl.Error, got %T", err)
	}
	if !deeplError.IsNotFound() {
		t.Errorf("expected IsNotFound, got code %d", deeplError.Code)
	}
}

func TestListGlossaryEntries(t *testing.T) {
	var captured *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello\thallo\nworld\tWelt\n"))
	}))
	t.Cleanup(server.Close)

	client := deepl.New("key", deepl.BaseURL(server.URL))
	entries, err := client.ListGlossaryEntries(context.Background(), "gl-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured.Header.Get("Accept") != "text/tab-separated-values" {
		t.Errorf("Accept = %q, want %q", captured.Header.Get("Accept"), "text/tab-separated-values")
	}
	if captured.Header.Get("Authorization") != "DeepL-Auth-Key key" {
		t.Errorf("Authorization = %q, want %q", captured.Header.Get("Authorization"), "DeepL-Auth-Key key")
	}

	want := []deepl.GlossaryEntry{
		{Source: "hello", Target: "hallo"},
		{Source: "world", Target: "Welt"},
	}
	if !reflect.DeepEqual(entries, want) {
		t.Errorf("entries = %+v, want %+v", entries, want)
	}
}

func TestListGlossaryEntries_empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := deepl.New("key", deepl.BaseURL(server.URL))
	entries, err := client.ListGlossaryEntries(context.Background(), "gl-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("entries = %+v, want nil", entries)
	}
}

func TestListGlossaryEntries_error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message":"not found"}`))
	}))
	t.Cleanup(server.Close)

	client := deepl.New("key", deepl.BaseURL(server.URL))
	_, err := client.ListGlossaryEntries(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var deeplError deepl.Error
	if !errors.As(err, &deeplError) {
		t.Fatalf("expected deepl.Error, got %T", err)
	}
	if !deeplError.IsNotFound() {
		t.Errorf("expected IsNotFound, got code %d", deeplError.Code)
	}
}

func TestListGlossaryEntries_malformedTSV(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"three columns", "a\tb\tc\n"},
		{"single column", "hello\n"},
		{"empty source", "\thallo\n"},
		{"empty target", "hello\t\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.body))
			}))
			t.Cleanup(server.Close)

			client := deepl.New("key", deepl.BaseURL(server.URL))
			_, err := client.ListGlossaryEntries(context.Background(), "gl-123")
			if err == nil {
				t.Fatal("expected error for malformed TSV, got nil")
			}
		})
	}
}

func TestDeleteGlossary(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "DELETE", "/gl-123", http.StatusNoContent, "", nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	err := client.DeleteGlossary(context.Background(), "gl-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteGlossary_notFound(t *testing.T) {
	server := newTestServer(t, glossaryHandler(t, "DELETE", "/nonexistent", http.StatusNotFound, `{"message":"not found"}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	err := client.DeleteGlossary(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var deeplError deepl.Error
	if !errors.As(err, &deeplError) {
		t.Fatalf("expected deepl.Error, got %T", err)
	}
	if !deeplError.IsNotFound() {
		t.Errorf("expected IsNotFound, got code %d", deeplError.Code)
	}
}

func TestTranslate_noTranslations_errorType(t *testing.T) {
	server := newTestServer(t, translateHandler(t, http.StatusOK, `{"translations": []}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	_, _, err := client.Translate(context.Background(), "text", deepl.German)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var deeplError deepl.Error
	if !errors.As(err, &deeplError) {
		t.Fatalf("expected deepl.Error, got %T: %v", err, err)
	}
}
