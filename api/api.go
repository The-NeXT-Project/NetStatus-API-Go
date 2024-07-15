package api

import (
	"encoding/json"
	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
	"net"
	"net/http"
	"time"
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

	status, msg := ping(request.URL.Query().Get("ip"), request.URL.Query().Get("port"))

	res, _ := json.Marshal(tcpingRes{
		Status:  status,
		Message: msg,
	})

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write(res)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func ping(ip string, port string) (status string, msg string) {
	timeout := time.Duration(int64(config.Config.TcpingTimeout) * int64(time.Millisecond))

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		return "false", "TCP connection failed"
	}

	if conn != nil {
		defer func(conn net.Conn) {
			_ = conn.Close()
		}(conn)
	}

	return "true", "TCP connection successful"
}
