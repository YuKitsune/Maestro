package spotify

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "spotify"

type Config struct {
	IsEnabled bool `mapstructure:"enabled"`
	ClientId string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

func (c *Config) Enabled() bool {
	return c.IsEnabled
}