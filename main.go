package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/sspanel-uim/NetStatus-API-Go/api"
	"github.com/sspanel-uim/NetStatus-API-Go/config"
	"strconv"
)

func main() {
	configFilePath := pflag.String("config", "", "config file path")
	pflag.Parse()

	if *configFilePath != "" {
		config.Config.AddConfigPath(*configFilePath)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Group("v1").GET("/tcping", api.Tcping)

	port := strconv.Itoa(config.Config.GetInt("port"))
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
