package model

type StreamingServiceKey string

func (s StreamingServiceKey) String() string {
	return string(s)
}

type StreamingService struct {
	Name    string
	Key     string
	Enabled bool
}
