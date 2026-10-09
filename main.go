package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/The-NeXT-Project/NetStatus-API-Go/api"
	"github.com/The-NeXT-Project/NetStatus-API-Go/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/spf13/pflag"
)

func main() {
	configFilePath := pflag.String("config", "", "config file, or directory containing config.json")
	pflag.Parse()

	if err := config.Load(*configFilePath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg := config.Config
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	if cfg.TrustProxy {
		// Forwarding headers can be set by anyone, trusting them without a proxy lets clients dodge the rate limit
		router.Use(middleware.RealIP)
	}
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	if cfg.RateLimit > 0 {
		router.Use(httprate.LimitByIP(cfg.RateLimit, time.Minute))
	}

	router.Route("/v1", func(router chi.Router) {
		router.Get("/tcping", api.TcpingV1)
	})

	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       time.Minute,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownDone := make(chan struct{})

	go func() {
		defer close(shutdownDone)
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("Listening on %s", server.Addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	// ListenAndServe returns as soon as Shutdown starts, wait for in-flight probes to finish
	<-shutdownDone
}
