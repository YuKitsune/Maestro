package db

type Config struct {
	Uri      string `mapstructure:"uri"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}
