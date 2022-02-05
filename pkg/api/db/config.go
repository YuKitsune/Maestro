package db

type Config struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"name"`
}
