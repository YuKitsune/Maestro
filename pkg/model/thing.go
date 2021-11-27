package model

type Thing interface {
	StreamingServiceThings() []StreamingServiceSpecificEntity
}