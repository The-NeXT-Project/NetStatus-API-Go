package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sspanel-uim/NetStatus-API-Go/api"
	"github.com/sspanel-uim/NetStatus-API-Go/config"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Group("v1").GET("/tcping", api.Tcping)
	port := strconv.Itoa(config.GetPort())

	err := r.Run(":" + port)
	if err != nil {
		return
	}
}
