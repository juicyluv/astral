package queue

import "github.com/spf13/viper"

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func NewConfig() *Config {
	return &Config{
		User:     viper.GetString("queue.user"),
		Password: viper.GetString("queue.password"),
		Host:     viper.GetString("queue.host"),
		Port:     viper.GetString("queue.port"),
		Name:     viper.GetString("queue.name"),
	}
}
