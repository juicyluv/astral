package configs

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// LoadConfigs loads .env file and config from given path.
// Returns an error if failed.
func LoadConfigs(configPath string) error {
	err := loadViperConfig(configPath)
	if err != nil {
		return err
	}

	return godotenv.Load()
}

// loadViperConfig tries to find config file and load it.
func loadViperConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	return viper.ReadInConfig()
}
