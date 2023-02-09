package config

import (
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/model"
)

type Deezer interface {
	Type() model.StreamingServiceType
	Name() string
	Enabled() bool
	LogoFileName() string
}

type deezerViperConfig struct {
	v *viper.Viper
}

func NewDeezerViperConfig(v *viper.Viper) Deezer {
	v.SetDefault("enabled", true)
	v.SetDefault("logo_file_name", "deezer.png")

	return &deezerViperConfig{v}
}

func (c *deezerViperConfig) Type() model.StreamingServiceType {
	return model.DeezerStreamingService
}

func (c *deezerViperConfig) Name() string {
	return "Deezer"
}

func (c *deezerViperConfig) Enabled() bool {
	return c.v.GetBool("enabled")
}

func (c *deezerViperConfig) LogoFileName() string {
	return c.v.GetString("logo_file_name")
}
