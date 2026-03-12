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

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func translateHandler(t *testing.T, statusCode int, body string, onRequest func(*http.Request)) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/translate" {
			t.Fatalf("expected /translate, got %s", r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %v", err)
		}
		if onRequest != nil {
			onRequest(r)
		}
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(body))
	}
}

func TestTranslate_requiredFields(t *testing.T) {
	var captured *http.Request
	server := newTestServer(t, translateHandler(t, http.StatusOK,
		`{"translations": [{"detected_source_language": "EN", "text": "Hallo"}]}`,
		func(r *http.Request) { captured = r },
	))

	client := deepl.New("an-auth-key", deepl.BaseURL(server.URL))
	_, _, err := client.Translate(context.Background(), "Hello", deepl.German)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured == nil {
		t.Fatal("request was not captured")
	}

	if got := captured.Header.Get("Authorization"); got != "DeepL-Auth-Key an-auth-key" {
		t.Errorf("Authorization = %q, want %q", got, "DeepL-Auth-Key an-auth-key")
	}
	if got := captured.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q, want %q", got, "application/x-www-form-urlencoded")
	}
	if got := captured.FormValue("target_lang"); got != string(deepl.German) {
		t.Errorf("target_lang = %q, want %q", got, string(deepl.German))
	}
	if got := captured.FormValue("text"); got != "Hello" {
		t.Errorf("text = %q, want %q", got, "Hello")
	}
}

func TestTranslate_options(t *testing.T) {
	tests := []struct {
		name  string
		opt   deepl.TranslateOption
		key   string
		value string
	}{
		{"SourceLang", deepl.SourceLang(deepl.English), "source_lang", string(deepl.English)},
		{"SplitSentences", deepl.SplitSentences(deepl.SplitNone), "split_sentences", deepl.SplitNone.Value()},
		{"PreserveFormatting", deepl.PreserveFormatting(true), "preserve_formatting", "1"},
		{"Formality", deepl.Formality(deepl.LessFormal), "formality", deepl.LessFormal.Value()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var captured *http.Request
			server := newTestServer(t, translateHandler(t, http.StatusOK,
				`{"translations": [{"detected_source_language": "EN", "text": "x"}]}`,
				func(r *http.Request) { captured = r },
			))

			client := deepl.New("key", deepl.BaseURL(server.URL))
			_, _, err := client.Translate(context.Background(), "text", deepl.German, tt.opt)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := captured.FormValue(tt.key); got != tt.value {
				t.Errorf("%s = %q, want %q", tt.key, got, tt.value)
			}
		})
	}
}

func TestTranslate_errors(t *testing.T) {
	codes := []int{
		http.StatusBadRequest,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestEntityTooLarge,
		http.StatusTooManyRequests,
		456,
		http.StatusServiceUnavailable,
	}

	for _, code := range codes {
		t.Run(http.StatusText(code), func(t *testing.T) {
			server := newTestServer(t, translateHandler(t, code, "{}", nil))

			client := deepl.New("key", deepl.BaseURL(server.URL))
			_, _, err := client.Translate(context.Background(), "text", deepl.German)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var deeplError deepl.Error
			if !errors.As(err, &deeplError) {
				t.Fatalf("expected deepl.Error, got %T", err)
			}
			if deeplError.Code != code {
				t.Errorf("Code = %d, want %d", deeplError.Code, code)
			}
		})
	}
}

func TestTranslate_noTranslations(t *testing.T) {
	server := newTestServer(t, translateHandler(t, http.StatusOK, `{"translations": []}`, nil))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	text, lang, err := client.Translate(context.Background(), "text", deepl.German)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if text != "" {
		t.Errorf("text = %q, want empty", text)
	}
	if lang != "" {
		t.Errorf("lang = %q, want empty", lang)
	}
}

func TestTranslate_withTranslation(t *testing.T) {
	server := newTestServer(t, translateHandler(t, http.StatusOK,
		`{"translations": [{"detected_source_language": "EN", "text": "Dies ist ein Beispieltext."}]}`,
		nil,
	))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	text, lang, err := client.Translate(context.Background(), "This is an example text.", deepl.German)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if text != "Dies ist ein Beispieltext." {
		t.Errorf("text = %q, want %q", text, "Dies ist ein Beispieltext.")
	}
	if lang != deepl.English {
		t.Errorf("lang = %q, want %q", lang, deepl.English)
	}
}

func TestTranslateMany_requiredFields(t *testing.T) {
	var captured *http.Request
	server := newTestServer(t, translateHandler(t, http.StatusOK, `{"translations": []}`,
		func(r *http.Request) { captured = r },
	))

	client := deepl.New("an-auth-key", deepl.BaseURL(server.URL))
	_, _ = client.TranslateMany(context.Background(), []string{"a", "b"}, deepl.German)
	if captured == nil {
		t.Fatal("request was not captured")
	}

	if got := captured.Header.Get("Authorization"); got != "DeepL-Auth-Key an-auth-key" {
		t.Errorf("Authorization = %q, want %q", got, "DeepL-Auth-Key an-auth-key")
	}
	if got := captured.Header.Get("Content-Type"); got != "application/x-www-form-urlencoded" {
		t.Errorf("Content-Type = %q, want %q", got, "application/x-www-form-urlencoded")
	}
	if got := captured.FormValue("target_lang"); got != string(deepl.German) {
		t.Errorf("target_lang = %q, want %q", got, string(deepl.German))
	}
	if got := captured.Form["text"]; !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("text = %v, want %v", got, []string{"a", "b"})
	}
}

func TestTranslateMany_withTranslations(t *testing.T) {
	server := newTestServer(t, translateHandler(t, http.StatusOK,
		`{"translations": [
			{"detected_source_language": "EN", "text": "Dies ist ein Beispiel."},
			{"detected_source_language": "FR", "text": "Dies ist ein anderer Text."}
		]}`,
		nil,
	))

	client := deepl.New("key", deepl.BaseURL(server.URL))
	translations, err := client.TranslateMany(
		context.Background(),
		[]string{"This is an example.", "C'est un autre texte."},
		deepl.German,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []deepl.Translation{
		{DetectedSourceLanguage: "EN", Text: "Dies ist ein Beispiel."},
		{DetectedSourceLanguage: "FR", Text: "Dies ist ein anderer Text."},
	}
	if !reflect.DeepEqual(translations, want) {
		t.Errorf("translations = %+v, want %+v", translations, want)
	}
}

// roundTripFunc adapts a function to http.RoundTripper for testing.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestTranslate_withCustomHTTPClient(t *testing.T) {
	called := false
	customClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			called = true
			return httptest.NewRecorder().Result(), nil
		}),
	}

	client := deepl.New("an-auth-key", deepl.HTTPClient(customClient))
	_, _, _ = client.Translate(context.Background(), "This is an example text.", deepl.German)

	if !called {
		t.Error("custom HTTP client was not called")
	}
}

func TestClient_HTTPClient(t *testing.T) {
	customClient := &http.Client{}
	client := deepl.New("an-auth-key", deepl.HTTPClient(customClient))
	if client.HTTPClient() != customClient {
		t.Error("HTTPClient() did not return the custom client")
	}
}

func TestClient_BaseURL(t *testing.T) {
	client := deepl.New("an-auth-key", deepl.BaseURL("base-url"))
	if got := client.BaseURL(); got != "base-url" {
		t.Errorf("BaseURL() = %q, want %q", got, "base-url")
	}
}

func TestClient_AuthKey(t *testing.T) {
	client := deepl.New("supersecure123")
	if got := client.AuthKey(); got != "supersecure123" {
		t.Errorf("AuthKey() = %q, want %q", got, "supersecure123")
	}
}

func TestClient_AuthKeyFromEnv(t *testing.T) {
	t.Setenv("DEEPL_AUTH_KEY", "env-key-123")
	client := deepl.New("")
	if got := client.AuthKey(); got != "env-key-123" {
		t.Errorf("AuthKey() = %q, want %q", got, "env-key-123")
	}
}

func TestError_classification(t *testing.T) {
	tests := []struct {
		code int
		name string
		check func(deepl.Error) bool
	}{
		{http.StatusBadRequest, "IsBadRequest", deepl.Error.IsBadRequest},
		{http.StatusForbidden, "IsUnauthorized", deepl.Error.IsUnauthorized},
		{http.StatusNotFound, "IsNotFound", deepl.Error.IsNotFound},
		{http.StatusRequestEntityTooLarge, "IsPayloadTooLarge", deepl.Error.IsPayloadTooLarge},
		{http.StatusTooManyRequests, "IsRateLimit", deepl.Error.IsRateLimit},
		{456, "IsQuotaExceeded", deepl.Error.IsQuotaExceeded},
		{http.StatusServiceUnavailable, "IsServiceUnavailable", deepl.Error.IsServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deepl.Error{Code: tt.code}
			if !tt.check(err) {
				t.Errorf("Error{Code: %d}.%s() = false, want true", tt.code, tt.name)
			}
		})
	}
}
