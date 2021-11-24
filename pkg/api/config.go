package api

type Config struct {
	Port     int    `mapstructure:"port"`
	CertFile string `mapstructure:"certificate-file"`
	KeyFile  string `mapstructure:"key-file"`
}
