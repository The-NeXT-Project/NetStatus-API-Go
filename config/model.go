package config

type ApiConfig struct {
	Port          int `mapstructure:"port"`
	ApiTimeout    int `mapstructure:"api_timeout"`
	TcpingTimeout int `mapstructure:"tcping_timeout"`
	RateLimit     int `mapstructure:"rate_limit"`
}
