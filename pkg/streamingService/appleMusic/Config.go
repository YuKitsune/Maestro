package appleMusic

import "github.com/yukitsune/maestro/pkg/model"

const Key model.StreamingServiceKey = "apple_music"

type Config struct {
	ServiceName         string `mapstructure:"name"`
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
	Token               string `mapstructure:"token"`
}

func (c *Config) Name() string {
	return c.ServiceName
}

func (c *Config) LogoFileName() string {
	return c.ServiceLogoFileName
}

func (c *Config) Enabled() bool {
	return c.IsEnabled
}
