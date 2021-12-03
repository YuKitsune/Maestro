package appleMusic

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "apple_music"

type Config struct {
	IsEnabled bool `mapstructure:"enabled"`
	Token string `mapstructure:"token"`
}

func (c *Config) Enabled() bool {
	return c.IsEnabled
}
