package deezer

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "deezer"

type Config struct {
	ServiceName         string `mapstructure:"name"`
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
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
