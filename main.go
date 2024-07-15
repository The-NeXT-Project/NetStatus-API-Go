package main

import (
	"context"
	"github.com/The-NeXT-Project/NetStatus-API-Go/api"
	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/spf13/pflag"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	configFilePath := pflag.String("config", "", "config file path")
	pflag.Parse()

	if *configFilePath != "" {
		config.Viper.AddConfigPath(*configFilePath)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(time.Duration(int64(config.Config.ApiTimeout) * int64(time.Second))))
	router.Use(httprate.LimitByIP(config.Config.RateLimit, 1*time.Minute))

	router.Route("/v1", func(router chi.Router) {
		router.Get("/tcping", api.TcpingV1)
	})

	port := ":" + strconv.Itoa(config.Config.Port)

	server := &http.Server{
		Addr:        port,
		BaseContext: func(net.Listener) context.Context { return context.Background() },
		Handler:     router,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
