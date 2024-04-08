package config

type ApiConfig struct {
	Port    int `mapstructure:"port"`
	Timeout int `mapstructure:"timeout"`
}
