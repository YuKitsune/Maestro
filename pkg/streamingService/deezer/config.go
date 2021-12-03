package deezer

import "maestro/pkg/model"

const Key model.StreamingServiceKey = "deezer"

type Config struct {
	IsEnabled bool `mapstructure:"enabled"`
}

func (c *Config) Enabled() bool {
	return c.IsEnabled
}