package spotify

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "spotify"

type Config struct {
	ServiceName string `mapstructure:"name"`
	ServiceArtworkLink string `mapstructure:"artwork_link"`
	IsEnabled bool `mapstructure:"enabled"`
	ClientId string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

func (c *Config) Name() string {
	return c.ServiceName
}

func (c *Config) ArtworkLink() string {
	return c.ServiceArtworkLink
}

func (c *Config) Enabled() bool {
	return c.IsEnabled
}