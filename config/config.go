package config

import "github.com/spf13/viper"

var config = viper.New()

func init() {
	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("/etc/netstatus-api-go/")
	config.AddConfigPath(".")

	config.SetDefault("port", 8080)
	config.SetDefault("timeout", 1000)

	err := config.ReadInConfig()
	if err != nil {
		return
	}
}

func GetPort() int {
	return config.GetInt("port")
}

func GetTimeout() int {
	return config.GetInt("timeout")
}
