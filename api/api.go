package api

import (
	"github.com/SSPanel-NeXT/NetStatus-API-Go/config"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

func Tcping(c *gin.Context) {
	if c.Query("ip") == "" {
		c.JSON(http.StatusBadRequest, tcpingRes{
			Status:  "false",
			Message: "Missing ip parameter",
		})

		return
	}

	if c.Query("port") == "" {
		c.JSON(http.StatusBadRequest, tcpingRes{
			Status:  "false",
			Message: "Missing port parameter",
		})

		return
	}

	status := "true"
	msg := ""
	status, msg = ping(c.Query("ip"), c.Query("port"))

	c.JSON(http.StatusOK, tcpingRes{
		Status:  status,
		Message: msg,
	})
}

func ping(ip string, port string) (status string, msg string) {
	timeout := time.Duration(int64(config.Config.GetInt("timeout")) * int64(time.Millisecond))

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
