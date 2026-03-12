package deepl_test

import (
	"testing"

	"github.com/solarhell/go-deepl"
)

func TestSplitSentence_Value_String(t *testing.T) {
	tests := map[deepl.SplitSentence]string{
		deepl.SplitNone:       "0",
		deepl.SplitDefault:    "1",
		deepl.SplitNoNewlines: "nonewlines",
	}

	for split, want := range tests {
		if got := split.Value(); got != want {
			t.Errorf("%v.Value() = %q, want %q", split, got, want)
		}
		if got := split.String(); got != want {
			t.Errorf("%v.String() = %q, want %q", split, got, want)
		}
	}
}
