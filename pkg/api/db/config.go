package db

type Config struct {
	Uri      string `mapstructure:"uri"`
	Database string `mapstructure:"name"`
}
