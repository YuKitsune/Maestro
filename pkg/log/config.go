package log

type Config struct {
	Level string `mapstructure:"level"`
	Loki  *struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"loki"`
}
