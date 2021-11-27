package model

type Thing interface {
	CollName() string
	GetLinks() Links
	SetLink(key StreamingServiceKey, link Link)
}
