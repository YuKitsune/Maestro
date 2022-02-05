package apiconfig

type Config struct {
	Port            int    `mapstructure:"port"`
	AssetsDirectory string `mapstructure:"assets_dir"`
	CertFile        string `mapstructure:"certificate-file"`
	KeyFile         string `mapstructure:"key-file"`
}
