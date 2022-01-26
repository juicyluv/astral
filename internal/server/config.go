package server

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

func NewConfig(configPath string) Config {
	return Config{
		Port:           ":" + viper.GetString("http.port"),
		ReadTimeout:    time.Second * time.Duration(viper.GetInt("http.readTimeout")),
		WriteTimeout:   time.Second * time.Duration(viper.GetInt("http.writeTimeout")),
		MaxHeaderBytes: viper.GetInt("http.maxHeaderBytes") << 20,
	}
}
