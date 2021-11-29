package model

import "testing"

func Test_NormalizeTrackName(t *testing.T) {
	values := []string {
		"My Track Name [Mixed]",
		"My Track Name (Mixed)",
		"My Track Name - Mixed",
		"My Track Name (feat. some dude) [Mixed]",
		"My Track Name (feat. some dude) (Mixed)",
		"My Track Name (feat. some dude) - Mixed",
	}

	exp := "My Track Name"

	norm := NewNameNormalizer()
	for _, value := range values {
		act := norm.NormalizeTrackName(value)
		if exp != act {
			t.Errorf("expected: \"%s\", found: \"%s\"", exp, act)
		}
	}
}

func Test_NormalizeAlbumName(t *testing.T) {
	values := []string {
		"My Track Name [Mixed] L.P.",
		"My Track Name (Mixed) - Single",
		"My Track Name - Mixed LP",
		"My Track Name (feat. some dude) [Mixed] - Single",
		"My Track Name (feat. some dude) (Mixed)",
		"My Track Name (feat. some dude) - Mixed - Single",
	}

	exp := "My Track Name"

	norm := NewNameNormalizer()
	for _, value := range values {
		act := norm.NormalizeAlbumName(value)
		if exp != act {
			t.Errorf("expected: \"%s\", found: \"%s\"", exp, act)
		}
	}
}