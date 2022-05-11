package apiconfig

type Config struct {
	Port            int    `mapstructure:"port"`
	AssetsDirectory string `mapstructure:"assets_dir"`
}
