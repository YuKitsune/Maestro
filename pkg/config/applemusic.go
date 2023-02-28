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
	v.SetDefault("services.apple_music.enabled", true)
	v.SetDefault("services.apple_music.logo_file_name", "apple_music.png")

	return &appleMusicViperConfig{v}
}

func (c *appleMusicViperConfig) Type() model.StreamingServiceType {
	return model.AppleMusicStreamingService
}

func (c *appleMusicViperConfig) Name() string {
	return "Apple Music"
}

func (c *appleMusicViperConfig) Enabled() bool {
	return c.v.GetBool("services.apple_music.enabled")
}

func (c *appleMusicViperConfig) LogoFileName() string {
	return c.v.GetString("services.apple_music.logo_file_name")
}

func (c *appleMusicViperConfig) Token() string {
	if !c.v.IsSet("services.apple_music.token") {
		panic("apple music token not set")
	}

	return c.v.GetString("services.apple_music.token")
}
