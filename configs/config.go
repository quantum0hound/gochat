package configs

import "github.com/spf13/viper"

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	// If a config file is found, read it in
	return viper.ReadInConfig()
}
