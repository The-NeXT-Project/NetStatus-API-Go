package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

var Config = &ApiConfig{}

// Load reads the config into Config. path may be a config file or a directory containing config.json;
// when empty, /etc/netstatus-api-go/ and the working directory are searched and defaults are used if neither has one.
func Load(path string) error {
	v := viper.New()
	v.SetConfigType("json")

	v.SetDefault("port", 8080)
	v.SetDefault("api_timeout", 3000)
	v.SetDefault("tcping_timeout", 1000)
	v.SetDefault("tcping_attempts", 1)
	v.SetDefault("rate_limit", 60)
	v.SetDefault("trust_proxy", false)
	v.SetDefault("allow_private", false)
	v.SetDefault("token", "")

	if path == "" {
		v.SetConfigName("config")
		v.AddConfigPath("/etc/netstatus-api-go/")
		v.AddConfigPath(".")
	} else if info, err := os.Stat(path); err == nil && info.IsDir() {
		v.SetConfigName("config")
		v.AddConfigPath(path)
	} else {
		v.SetConfigFile(path)
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if path != "" || !errors.As(err, &notFound) {
			return err
		}

		log.Print("No config file found, using defaults")
	}

	cfg := &ApiConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	if err := cfg.validate(); err != nil {
		return err
	}

	Config = cfg

	return nil
}

func (c *ApiConfig) validate() error {
	switch {
	case c.Port < 1 || c.Port > 65535:
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	case c.ApiTimeout <= 0:
		return fmt.Errorf("api_timeout must be positive, got %d", c.ApiTimeout)
	case c.TcpingTimeout <= 0:
		return fmt.Errorf("tcping_timeout must be positive, got %d", c.TcpingTimeout)
	case c.TcpingAttempts < 1:
		return fmt.Errorf("tcping_attempts must be at least 1, got %d", c.TcpingAttempts)
	case c.TcpingAttempts*c.TcpingTimeout > c.ApiTimeout:
		// Otherwise the last attempts get cut off by api_timeout and the request ends without a verdict
		return fmt.Errorf("tcping_attempts * tcping_timeout (%d ms) must not exceed api_timeout (%d ms)",
			c.TcpingAttempts*c.TcpingTimeout, c.ApiTimeout)
	case c.RateLimit < 0:
		return fmt.Errorf("rate_limit must not be negative, got %d", c.RateLimit)
	}

	return nil
}
