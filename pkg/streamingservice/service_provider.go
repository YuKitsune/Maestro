package streamingservice

import "github.com/yukitsune/maestro/pkg/model"

type ServiceProvider interface {
	GetService(key model.StreamingServiceKey) StreamingService
	ListServices() []StreamingService
	GetConfig(key model.StreamingServiceKey) ServiceConfig
	ListConfigs() Config
}
