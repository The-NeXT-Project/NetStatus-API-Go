package api

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"time"

	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
)

// Ranges that IsGlobalUnicast and IsPrivate don't cover but are still not public internet hosts
var nonPublicPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
}

// TcpingV1 answers whether ip:port accepts TCP connections from here.
//
// NeXT-Panel only reads the body of 2xx responses and treats anything else as "no verdict", so a
// 200 with status "false" must only ever mean the target was actually unreachable.
func TcpingV1(writer http.ResponseWriter, request *http.Request) {
	cfg := config.Config
	query := request.URL.Query()

	if cfg.Token != "" && subtle.ConstantTimeCompare([]byte(query.Get("token")), []byte(cfg.Token)) != 1 {
		writeError(writer, http.StatusUnauthorized, "Invalid token")
		return
	}

	ipStr := query.Get("ip")
	if ipStr == "" {
		writeError(writer, http.StatusBadRequest, "Missing ip parameter")
		return
	}

	portStr := query.Get("port")
	if portStr == "" {
		writeError(writer, http.StatusBadRequest, "Missing port parameter")
		return
	}

	addr, err := netip.ParseAddr(ipStr)
	if err != nil || addr.Zone() != "" {
		writeError(writer, http.StatusBadRequest, "Invalid IP address format")
		return
	}

	addr = addr.Unmap()
	if !cfg.AllowPrivate && !isPublic(addr) {
		writeError(writer, http.StatusForbidden, "Target address is not allowed")
		return
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil || port == 0 {
		writeError(writer, http.StatusBadRequest, "Invalid port number")
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), time.Duration(cfg.ApiTimeout)*time.Millisecond)
	defer cancel()

	latency, ok := ping(ctx, netip.AddrPortFrom(addr, uint16(port)))

	switch {
	case ok:
		ms := latency.Milliseconds()
		writeJSON(writer, http.StatusOK, tcpingRes{
			Status:  "true",
			Time:    &ms,
			Message: "TCP connection successful",
		})
	case ctx.Err() != nil:
		// Cut off by api_timeout or the client going away, not a measured failure
		writeError(writer, http.StatusGatewayTimeout, "Probe did not finish in time")
	default:
		writeJSON(writer, http.StatusOK, tcpingRes{
			Status:  "false",
			Message: "TCP connection failed",
		})
	}
}

// ping makes up to tcping_attempts connection attempts and returns the latency of the first one that succeeds.
func ping(ctx context.Context, target netip.AddrPort) (time.Duration, bool) {
	dialer := net.Dialer{Timeout: time.Duration(config.Config.TcpingTimeout) * time.Millisecond}

	for range config.Config.TcpingAttempts {
		startTime := time.Now()

		conn, err := dialer.DialContext(ctx, "tcp", target.String())
		if err == nil {
			latency := time.Since(startTime)
			_ = conn.Close()

			return latency, true
		}

		if ctx.Err() != nil {
			break
		}
	}

	return 0, false
}

func isPublic(addr netip.Addr) bool {
	if !addr.IsGlobalUnicast() || addr.IsPrivate() {
		return false
	}

	for _, prefix := range nonPublicPrefixes {
		if prefix.Contains(addr) {
			return false
		}
	}

	return true
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, tcpingRes{
		Status:  "error",
		Message: message,
	})
}

func writeJSON(writer http.ResponseWriter, status int, res tcpingRes) {
	body, _ := json.Marshal(res)

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(status)
	// Headers are already sent, nothing useful can be done if the write fails
	_, _ = writer.Write(body)
}
