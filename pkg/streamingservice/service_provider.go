package streamingservice

import (
	"github.com/yukitsune/maestro/pkg/config"
	"github.com/yukitsune/maestro/pkg/model"
)

type ServiceProvider interface {
	GetService(model.StreamingServiceType) (StreamingService, error)
	ListServices() (StreamingServices, error)
	GetConfig(model.StreamingServiceType) (config.Service, error)
	ListConfigs() map[model.StreamingServiceType]config.Service
}
