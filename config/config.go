package config

import "github.com/spf13/viper"

var (
	Viper  = viper.New()
	Config = &ApiConfig{}
)

func init() {
	Viper.SetConfigName("config")
	Viper.SetConfigType("json")
	Viper.AddConfigPath("/etc/netstatus-api-go/")
	Viper.AddConfigPath(".")

	Viper.SetDefault("port", 8080)
	Viper.SetDefault("api_timeout", 3000)
	Viper.SetDefault("tcping_timeout", 1000)
	Viper.SetDefault("rate_limit", 60)

	err := Viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = Viper.Unmarshal(&Config)
	if err != nil {
		panic(err)
	}
}
