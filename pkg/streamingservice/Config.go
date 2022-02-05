package streamingservice

import (
	"github.com/yukitsune/maestro/pkg/model"
)

type Config map[model.StreamingServiceKey]ServiceConfig

type ServiceConfig interface {
	Name() string
	LogoFileName() string
	Enabled() bool
}
