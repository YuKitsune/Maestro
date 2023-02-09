package config

import (
	"github.com/spf13/viper"
	"github.com/yukitsune/maestro/pkg/model"
)

type AppleMusic interface {
	Type() model.StreamingServiceType
	Name() string
	Enabled() bool
	LogoFileName() string
	Token() string
}

type appleMusicViperConfig struct {
	v *viper.Viper
}

func NewAppleMusicViperConfig(v *viper.Viper) AppleMusic {
	v.SetDefault("enabled", true)
	v.SetDefault("logo_file_name", "deezer.png")

	return &appleMusicViperConfig{v}
}

func (c *appleMusicViperConfig) Type() model.StreamingServiceType {
	return model.AppleMusicStreamingService
}

func (c *appleMusicViperConfig) Name() string {
	return "Apple Music"
}

func (c *appleMusicViperConfig) Enabled() bool {
	return c.v.GetBool("enabled")
}

func (c *appleMusicViperConfig) LogoFileName() string {
	return c.v.GetString("logo_file_name")
}

func (c *appleMusicViperConfig) Token() string {
	if !c.v.IsSet("token") {
		panic("apple music token not set")
	}

	return c.v.GetString("token")
}
