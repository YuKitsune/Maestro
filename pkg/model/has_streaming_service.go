package model

type HasStreamingService interface {
	GetSource() StreamingServiceType
}
