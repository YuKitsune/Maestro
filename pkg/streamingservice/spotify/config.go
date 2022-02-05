package spotify

import "github.com/yukitsune/maestro/pkg/model"

const Key model.StreamingServiceKey = "spotify"

type Config struct {
	ServiceName         string `mapstructure:"name"`
	ServiceLogoFileName string `mapstructure:"logo_file_name"`
	IsEnabled           bool   `mapstructure:"enabled"`
	ClientID            string `mapstructure:"client_id"`
	ClientSecret        string `mapstructure:"client_secret"`
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
