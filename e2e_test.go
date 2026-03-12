package deepl_test

import (
	"context"
	"os"
	"testing"

	"github.com/solarhell/go-deepl"
)

func TestTranslate_withoutSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	translated, sourceLang, err := client.Translate(
		context.Background(),
		"This is an example text.",
		deepl.German,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if translated != "Dies ist ein Beispieltext." {
		t.Errorf("translated = %q, want %q", translated, "Dies ist ein Beispieltext.")
	}
	if sourceLang != deepl.English {
		t.Errorf("sourceLang = %q, want %q", sourceLang, deepl.English)
	}
}

func TestTranslate_showBilledCharacters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	translations, err := client.TranslateMany(
		context.Background(),
		[]string{"This is an example text."},
		deepl.German,
		deepl.ShowBilledChars(true),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(translations) != 1 {
		t.Fatalf("got %d translations, want 1", len(translations))
	}
	if translations[0].Text != "Dies ist ein Beispieltext." {
		t.Errorf("text = %q, want %q", translations[0].Text, "Dies ist ein Beispieltext.")
	}
	if deepl.Language(translations[0].DetectedSourceLanguage) != deepl.English {
		t.Errorf("detected language = %q, want %q", translations[0].DetectedSourceLanguage, deepl.English)
	}
	if translations[0].BilledCharacters <= 0 {
		t.Errorf("billed characters = %d, want > 0", translations[0].BilledCharacters)
	}
}

func TestTranslate_withSourceLang(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	_, sourceLang, err := client.Translate(
		context.Background(),
		"Voici un exemple de texte.",
		deepl.German,
		deepl.SourceLang(deepl.English),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sourceLang != deepl.English {
		t.Errorf("sourceLang = %q, want %q", sourceLang, deepl.English)
	}

	// we don't validate the translated text, because the translation behaviour
	// for an invalid source language is not defined
}

func TestHTMLTagHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test.")
	}

	client := deepl.New(getAuthKey(t), getOpts()...)

	res, _, err := client.Translate(
		context.Background(),
		`<p alt="This is a test.">This is a test.</p>`,
		deepl.German,
		deepl.TagHandling(deepl.HTMLTagHandling),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `<p alt="This is a test.">Dies ist ein Test.</p>`
	if res != want {
		t.Errorf("result = %q, want %q", res, want)
	}
}

func getOpts(opts ...deepl.ClientOption) []deepl.ClientOption {
	apiEndpoint := os.Getenv("DEEPL_API_ENDPOINT")
	ret := opts
	if apiEndpoint != "" {
		ret = append(ret, deepl.BaseURL(apiEndpoint))
	}
	return ret
}

func getAuthKey(t *testing.T) string {
	authKey := os.Getenv("DEEPL_AUTH_KEY")
	if authKey == "" {
		t.Fatal("Set the DEEPL_AUTH_KEY environment variable before running the integration tests.")
	}
	return authKey
}
