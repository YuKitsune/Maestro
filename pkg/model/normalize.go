package model

import (
	"regexp"
	"strings"
)

type NameNormalizer struct {
	featRegex *regexp.Regexp
	variationRegex *regexp.Regexp
	albumTypeRegex *regexp.Regexp
}

func NewNameNormalizer() *NameNormalizer {

	return &NameNormalizer{
		albumTypeRegex: regexp.MustCompile("\\s(([EeLl]\\.?[Pp].?)|(-\\s[Ss]ingle))"),
		variationRegex: regexp.MustCompile("(\\s+((\\[[\\w\\s]+\\])|(\\([\\w\\s]+\\))|(-[\\w\\s]+)))$"),
		featRegex:      regexp.MustCompile("(\\([Ff]eat(\\.|uring).+\\))"),
	}
}

func (n *NameNormalizer) NormalizeAlbumName(name string) string {

	patterns := []*regexp.Regexp {
		n.albumTypeRegex,
		n.variationRegex,
		n.featRegex,
	}

	for _, pattern := range patterns {
		index := pattern.FindStringIndex(name)
		if len(index) > 0 {
			start := index[0]
			end := index[1]

			name = name[0:start] + name[end:len(name)]
		}
	}

	name = strings.Trim(name," ")

	return name
}

func (n *NameNormalizer) NormalizeTrackName(name string) string {

	patterns := []*regexp.Regexp {
		n.variationRegex,
		n.featRegex,
	}

	for _, pattern := range patterns {
		index := pattern.FindStringIndex(name)
		if len(index) > 0 {
			start := index[0]
			end := index[1]

			name = name[0:start] + name[end:len(name)]
		}
	}

	name = strings.Trim(name," ")

	return name
}