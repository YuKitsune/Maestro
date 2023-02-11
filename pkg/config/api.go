package config

import "github.com/spf13/viper"

type API interface {
	Port() int
	AssetsDirectory() string
}

type apiViperConfig struct {
	v *viper.Viper
}

func NewApiViperConfig(v *viper.Viper) API {
	v.SetDefault("api.port", 8182)
	v.SetDefault("api.assets_dir", "/assets")

	return &apiViperConfig{v}
}

func (c *apiViperConfig) Port() int {
	return c.v.GetInt("api.port")
}

func (c *apiViperConfig) AssetsDirectory() string {
	return c.v.GetString("api.assets_dir")
}
