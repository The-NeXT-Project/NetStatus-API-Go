package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
)

func TcpingV1(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Query().Get("ip") == "" {
		res, _ := json.Marshal(tcpingRes{
			Status:  "false",
			Message: "Missing ip parameter",
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write(res)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	if request.URL.Query().Get("port") == "" {
		res, _ := json.Marshal(tcpingRes{
			Status:  "false",
			Message: "Missing port parameter",
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write(res)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	ip := request.URL.Query().Get("ip")
	portStr := request.URL.Query().Get("port")

	// Sanitize IP address
	if net.ParseIP(ip) == nil {
		res, _ := json.Marshal(tcpingRes{
			Status: "false",
			Message: "Invalid IP address format",
		})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write(res)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Sanitize port
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		res, _ := json.Marshal(tcpingRes{
			Status: "false",
			Message: "Invalid port number",
		})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write(res) // Error handling for Write is already present below
		return
	}

	status, latency, msg := ping(ip, portStr)

	res, _ := json.Marshal(tcpingRes{
		Status:  status,
		Time:    latency,
		Message: msg,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(res)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func ping(ip string, port string) (string, int, string) {
	timeout := time.Duration(int64(config.Config.TcpingTimeout) * int64(time.Millisecond))
	startTime := time.Now()

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return "false", config.Config.TcpingTimeout, "TCP connection failed"
	}

	if conn != nil {
		defer func(conn net.Conn) {
			_ = conn.Close()
		}(conn)
	}

	return "true", int(time.Since(startTime).Milliseconds()), "TCP connection successful"
}
