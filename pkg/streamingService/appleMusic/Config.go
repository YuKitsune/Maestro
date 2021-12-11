package appleMusic

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "apple_music"

type Config struct {
	ServiceName        string `mapstructure:"name"`
	ServiceArtworkLink string `mapstructure:"artwork_link"`
	IsEnabled          bool   `mapstructure:"enabled"`
	Token              string `mapstructure:"token"`
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
