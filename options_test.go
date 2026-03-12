package deepl_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/solarhell/go-deepl"
)

func TestSourceLang(t *testing.T) {
	vals := make(url.Values)
	if got := vals.Get("source_lang"); got != "" {
		t.Errorf("initial source_lang = %q, want empty", got)
	}
	deepl.SourceLang(deepl.German)(vals)
	if got := vals.Get("source_lang"); got != string(deepl.German) {
		t.Errorf("source_lang = %q, want %q", got, string(deepl.German))
	}
}

func TestShowBilledChars(t *testing.T) {
	vals := make(url.Values)
	deepl.ShowBilledChars(true)(vals)
	if got := vals.Get("show_billed_characters"); got != "1" {
		t.Errorf("show_billed_characters = %q, want %q", got, "1")
	}
	deepl.ShowBilledChars(false)(vals)
	if got := vals.Get("show_billed_characters"); got != "0" {
		t.Errorf("show_billed_characters = %q, want %q", got, "0")
	}
}

func TestSplitSentences(t *testing.T) {
	splits := []deepl.SplitSentence{
		deepl.SplitNone,
		deepl.SplitDefault,
		deepl.SplitNoNewlines,
	}

	for _, split := range splits {
		t.Run(split.String(), func(t *testing.T) {
			vals := make(url.Values)
			deepl.SplitSentences(split)(vals)
			if got := vals.Get("split_sentences"); got != split.Value() {
				t.Errorf("split_sentences = %q, want %q", got, split.Value())
			}
		})
	}
}

func TestPreserveFormatting(t *testing.T) {
	vals := make(url.Values)
	deepl.PreserveFormatting(true)(vals)
	if got := vals.Get("preserve_formatting"); got != "1" {
		t.Errorf("preserve_formatting = %q, want %q", got, "1")
	}
	deepl.PreserveFormatting(false)(vals)
	if got := vals.Get("preserve_formatting"); got != "0" {
		t.Errorf("preserve_formatting = %q, want %q", got, "0")
	}
}

func TestFormality(t *testing.T) {
	formalities := []deepl.Formal{
		deepl.DefaultFormal,
		deepl.LessFormal,
		deepl.MoreFormal,
	}

	for _, f := range formalities {
		t.Run(f.String(), func(t *testing.T) {
			vals := make(url.Values)
			deepl.Formality(f)(vals)
			if got := vals.Get("formality"); got != f.Value() {
				t.Errorf("formality = %q, want %q", got, f.Value())
			}
		})
	}
}

func TestTagHandling(t *testing.T) {
	strategies := []deepl.TagHandlingStrategy{
		deepl.DefaultTagHandling,
		deepl.XMLTagHandling,
	}

	for _, s := range strategies {
		t.Run(s.String(), func(t *testing.T) {
			vals := make(url.Values)
			deepl.TagHandling(s)(vals)
			if got := vals.Get("tag_handling"); got != s.Value() {
				t.Errorf("tag_handling = %q, want %q", got, s.Value())
			}
		})
	}
}

func TestIgnoreTags(t *testing.T) {
	tags := []string{"foo", "bar", "baz"}
	vals := make(url.Values)
	deepl.IgnoreTags(tags...)(vals)

	want := strings.Join(tags, ",")
	if got := vals.Get("ignore_tags"); got != want {
		t.Errorf("ignore_tags = %q, want %q", got, want)
	}
}
