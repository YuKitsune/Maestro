package model

type StreamingServiceKey string

func (s StreamingServiceKey) String() string {
	return string(s)
}
