package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
)

func setConfig(t *testing.T, cfg config.ApiConfig) {
	t.Helper()

	old := config.Config
	config.Config = &cfg
	t.Cleanup(func() { config.Config = old })
}

func defaultConfig() config.ApiConfig {
	return config.ApiConfig{ApiTimeout: 3000, TcpingTimeout: 1000, TcpingAttempts: 1}
}

func do(t *testing.T, ctx context.Context, query string) (int, tcpingRes) {
	t.Helper()

	request := httptest.NewRequestWithContext(ctx, http.MethodGet, "/v1/tcping?"+query, nil)
	recorder := httptest.NewRecorder()
	TcpingV1(recorder, request)

	var res tcpingRes
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not JSON: %q", recorder.Body.String())
	}

	return recorder.Code, res
}

func TestTcpingV1Rejects(t *testing.T) {
	setConfig(t, defaultConfig())

	cases := []struct {
		name   string
		query  string
		status int
	}{
		{"missing ip", "port=443", http.StatusBadRequest},
		{"missing port", "ip=1.1.1.1", http.StatusBadRequest},
		{"hostname", "ip=example.com&port=443", http.StatusBadRequest},
		{"zoned ipv6", "ip=fe80::1%25eth0&port=443", http.StatusBadRequest},
		{"port zero", "ip=1.1.1.1&port=0", http.StatusBadRequest},
		{"port too large", "ip=1.1.1.1&port=65536", http.StatusBadRequest},
		{"port not a number", "ip=1.1.1.1&port=https", http.StatusBadRequest},
		{"loopback", "ip=127.0.0.1&port=443", http.StatusForbidden},
		{"private", "ip=192.168.1.1&port=443", http.StatusForbidden},
		{"cgnat", "ip=100.64.0.1&port=443", http.StatusForbidden},
		{"link local metadata", "ip=169.254.169.254&port=80", http.StatusForbidden},
		{"mapped loopback", "ip=::ffff:127.0.0.1&port=443", http.StatusForbidden},
		{"ipv6 loopback", "ip=::1&port=443", http.StatusForbidden},
		{"unspecified", "ip=0.0.0.0&port=443", http.StatusForbidden},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, res := do(t, context.Background(), c.query)
			if code != c.status || res.Status != "error" {
				t.Errorf("got %d %q, want %d \"error\"", code, res.Status, c.status)
			}
		})
	}
}

func TestTcpingV1Token(t *testing.T) {
	cfg := defaultConfig()
	cfg.Token = "secret"
	setConfig(t, cfg)

	for _, query := range []string{"ip=1.1.1.1&port=443", "ip=1.1.1.1&port=443&token=wrong"} {
		if code, _ := do(t, context.Background(), query); code != http.StatusUnauthorized {
			t.Errorf("%s: got %d, want 401", query, code)
		}
	}

	// A correct token gets past auth to the address check
	if code, _ := do(t, context.Background(), "ip=127.0.0.1&port=443&token=secret"); code != http.StatusForbidden {
		t.Errorf("got %d, want 403", code)
	}
}

func TestTcpingV1Probe(t *testing.T) {
	cfg := defaultConfig()
	cfg.AllowPrivate = true
	setConfig(t, cfg)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	openPort := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)

	closed, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	closedPort := strconv.Itoa(closed.Addr().(*net.TCPAddr).Port)
	_ = closed.Close()

	t.Run("reachable", func(t *testing.T) {
		code, res := do(t, context.Background(), "ip=127.0.0.1&port="+openPort)
		if code != http.StatusOK || res.Status != "true" || res.Time == nil {
			t.Errorf("got %d %+v, want 200 \"true\" with time", code, res)
		}
	})

	t.Run("unreachable", func(t *testing.T) {
		code, res := do(t, context.Background(), "ip=127.0.0.1&port="+closedPort)
		if code != http.StatusOK || res.Status != "false" || res.Time != nil {
			t.Errorf("got %d %+v, want 200 \"false\" without time", code, res)
		}
	})

	t.Run("cut off is not a verdict", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		code, res := do(t, ctx, "ip=127.0.0.1&port="+openPort)
		if code != http.StatusGatewayTimeout || res.Status != "error" {
			t.Errorf("got %d %+v, want 504 \"error\"", code, res)
		}
	})
}
