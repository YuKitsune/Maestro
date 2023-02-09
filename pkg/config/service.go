package config

import "github.com/yukitsune/maestro/pkg/model"

type Service interface {
	Type() model.StreamingServiceType
	Name() string
	LogoFileName() string
	Enabled() bool
}
