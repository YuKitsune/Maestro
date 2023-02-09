package config

import "github.com/spf13/viper"

type Database interface {
	Uri() string
	Name() string
}

type databaseViperConfig struct {
	v *viper.Viper
}

func NewDatabaseViperConfig(v *viper.Viper) Database {
	return &databaseViperConfig{v}
}

func (c *databaseViperConfig) Uri() string {
	if !c.v.IsSet("uri") {
		panic("database uri not set")
	}

	return c.v.GetString("uri")
}

func (c *databaseViperConfig) Name() string {
	if !c.v.IsSet("name") {
		panic("database name not set")
	}

	return c.v.GetString("name")
}
