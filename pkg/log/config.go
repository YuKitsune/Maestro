package log

type Config struct {
	Level            string `mapstructure:"level"`
	UseJsonFormatter bool   `mapstructure:"use_json_format"`
}
