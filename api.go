package main

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

type tcpingRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var port = "8080"
var timeout = 1 * time.Second

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Group("v1").GET("/tcping", tcping)

	err := r.Run(":" + port)

	if err != nil {
		return
	}
}

func tcping(c *gin.Context) {
	if c.Query("ip") == "" {
		c.JSON(http.StatusOK, tcpingRes{
			Status:  "error",
			Message: "Missing ip parameter",
		})

		return
	}

	if c.Query("port") == "" {
		c.JSON(http.StatusOK, tcpingRes{
			Status:  "error",
			Message: "Missing ip parameter",
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
