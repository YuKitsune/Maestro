package model

type HasStreamingService interface {
	GetSource() StreamingServiceKey
}
