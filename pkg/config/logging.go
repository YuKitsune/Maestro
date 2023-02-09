package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
)

type Logging interface {
	Level() logrus.Level
	Loki() Loki
}

type Loki interface {
	Enabled() bool
	Host() string
	Labels() map[string]string
}

type loggingViperConfig struct {
	v    *viper.Viper
	loki Loki
}

func NewLoggingViperConfig(v *viper.Viper) Logging {
	v.SetDefault("level", "info")
	return &loggingViperConfig{
		v,
		NewLokiViperConfig(v.Sub("loki")),
	}
}

func (c *loggingViperConfig) Level() logrus.Level {

	level := c.v.GetString("level")
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		log.Fatalf("unable to parse log level: %s", err)
	}

	return logrusLevel
}

func (c *loggingViperConfig) Loki() Loki {
	return c.loki
}

type lokiViperConfig struct {
	v *viper.Viper
}

func NewLokiViperConfig(v *viper.Viper) Loki {
	return &lokiViperConfig{v}
}

func (c *lokiViperConfig) Enabled() bool {
	if c.v == nil {
		return false
	}

	return c.v.IsSet("host")
}

func (c *lokiViperConfig) Host() string {
	if !c.Enabled() {
		panic("loki not enabled")
	}

	if !c.v.IsSet("host") {
		panic("loki host not set")
	}

	return c.v.GetString("host")
}

func (c *lokiViperConfig) Labels() map[string]string {
	if !c.Enabled() {
		panic("loki not enabled")
	}

	return c.v.GetStringMapString("labels")
}
