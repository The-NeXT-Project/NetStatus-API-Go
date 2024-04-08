package config

import "github.com/spf13/viper"

var (
	Config    = viper.New()
	apiConfig = &ApiConfig{}
)

func init() {
	Config.SetConfigName("Config")
	Config.SetConfigType("json")
	Config.AddConfigPath("/etc/netstatus-api-go/")
	Config.AddConfigPath(".")

	Config.SetDefault("port", 8080)
	Config.SetDefault("timeout", 1000)

	err := Config.Unmarshal(&apiConfig)
	if err != nil {
		return
	}
}
