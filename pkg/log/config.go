package log

type Config struct {
	Level string `mapstructure:"level"`
	Loki  *struct {
		Host   string            `mapstructure:"host"`
		Labels map[string]string `mapstructure:"labels"`
	} `mapstructure:"loki"`
}
