package levenshtein

import (
	"testing"
)

func Test_Levenshtein(t *testing.T) {
	orig := "My Awesome Song! [Mixed]"
	in := "My Awesome Song! - Mixed"

	d := Levenshtein(orig, in)
	if d != 3 {
		t.Fail()
	}
}

func Test_Levenshtein_CompletelyDifferent(t *testing.T) {
	orig := "My Awesome Song! [Mixed]"
	in := "A completely different thing"

	d := Levenshtein(orig, in)
	if d != 26 {
		t.Fail()
	}
}
