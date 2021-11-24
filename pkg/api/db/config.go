package db

import "fmt"

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func (c *Config) ConnectionString() string {
	// Todo: Hand this to a provider so we can use different databases
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", c.Host, c.Port, c.User, c.Password, c.Database)
}
