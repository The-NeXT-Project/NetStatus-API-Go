package config

type ApiConfig struct {
	Port int `mapstructure:"port"`
	// ApiTimeout bounds a whole /v1/tcping request, in milliseconds
	ApiTimeout int `mapstructure:"api_timeout"`
	// TcpingTimeout bounds a single TCP connection attempt, in milliseconds
	TcpingTimeout int `mapstructure:"tcping_timeout"`
	// TcpingAttempts is how many connection attempts are made before reporting the target as unreachable
	TcpingAttempts int `mapstructure:"tcping_attempts"`
	// RateLimit is requests per minute per client IP, 0 disables it
	RateLimit int `mapstructure:"rate_limit"`
	// TrustProxy takes the client IP from X-Forwarded-For / X-Real-IP, only enable it behind a reverse proxy
	TrustProxy bool `mapstructure:"trust_proxy"`
	// AllowPrivate allows probing loopback, private and other non-public addresses
	AllowPrivate bool `mapstructure:"allow_private"`
	// Token, when set, must be passed as the token query parameter
	Token string `mapstructure:"token"`
}
