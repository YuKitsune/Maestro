package streamingService

import (
	"maestro/pkg/model"
)

type Config map[model.StreamingServiceKey]ServiceConfig

type ServiceConfig interface {
	Enabled() bool
}
