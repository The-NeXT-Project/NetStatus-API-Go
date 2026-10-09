package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return dir
}

func TestLoad(t *testing.T) {
	dir := writeConfig(t, `{"port": 9090, "token": "secret"}`)

	for _, path := range []string{dir, filepath.Join(dir, "config.json")} {
		if err := Load(path); err != nil {
			t.Fatalf("Load(%q): %v", path, err)
		}

		if Config.Port != 9090 || Config.Token != "secret" || Config.ApiTimeout != 3000 || Config.TcpingAttempts != 1 {
			t.Errorf("Load(%q) = %+v", path, *Config)
		}
	}
}

func TestLoadMissingExplicitPath(t *testing.T) {
	if err := Load(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Error("expected an error for a missing config file")
	}

	if err := Load(t.TempDir()); err == nil {
		t.Error("expected an error for a directory without config.json")
	}
}

func TestLoadValidation(t *testing.T) {
	cases := map[string]string{
		`{"port": 0}`:                                   "port",
		`{"api_timeout": 0}`:                            "api_timeout",
		`{"tcping_attempts": 0}`:                        "tcping_attempts",
		`{"tcping_attempts": 4}`:                        "must not exceed api_timeout",
		`{"rate_limit": -1}`:                            "rate_limit",
		`{"tcping_timeout": 2000, "api_timeout": 1000}`: "must not exceed api_timeout",
	}

	for content, want := range cases {
		err := Load(writeConfig(t, content))
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want error mentioning %q", content, err, want)
		}
	}
}
