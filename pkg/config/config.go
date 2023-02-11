package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config interface {
	API() API
	Database() Database
	Logging() Logging
	Services() Services
	Debug() string
}

type viperConfig struct {
	v        *viper.Viper
	api      API
	database Database
	logging  Logging
	services Services
}

func NewViperConfig(v *viper.Viper) Config {
	return &viperConfig{
		v: v,
		// Todo: Update this to use sub once viper bug is fixed
		api:      NewApiViperConfig(v),
		database: NewDatabaseViperConfig(v),
		logging:  NewLoggingViperConfig(v),
		services: NewServicesViperConfig(v)}
}

func (c *viperConfig) API() API {
	return c.api
}

func (c *viperConfig) Database() Database {
	return c.database
}

func (c *viperConfig) Logging() Logging {
	return c.logging
}

func (c *viperConfig) Services() Services {
	return c.services
}

func (c *viperConfig) Debug() string {
	return fmt.Sprintf("%#v", c.v.AllSettings())
}
