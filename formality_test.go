package deepl_test

import (
	"testing"

	"github.com/solarhell/go-deepl"
)

func TestFormality_Value_String(t *testing.T) {
	tests := map[deepl.Formal]string{
		deepl.DefaultFormal: "default",
		deepl.LessFormal:    "less",
		deepl.MoreFormal:    "more",
	}

	for f, want := range tests {
		if got := f.Value(); got != want {
			t.Errorf("%v.Value() = %q, want %q", f, got, want)
		}
		if got := f.String(); got != want {
			t.Errorf("%v.String() = %q, want %q", f, got, want)
		}
	}
}
