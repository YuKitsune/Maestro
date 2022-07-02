package streamingservice

import "github.com/yukitsune/maestro/pkg/model"

type ServiceProvider interface {
	GetService(key model.StreamingServiceKey) (StreamingService, error)
	ListServices() (StreamingServices, error)
	GetConfig(key model.StreamingServiceKey) (ServiceConfig, error)
	ListConfigs() Config
}
