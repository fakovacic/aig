package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	"github.com/fakovacic/aig/internal/app"
	"github.com/fakovacic/aig/internal/health"
	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/internal/middlewares"
	"github.com/joho/godotenv"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttprouter"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	var (
		httpAddr   = "0.0.0.0:" + os.Getenv("APP_PORT")
		healthAddr = "0.0.0.0:" + os.Getenv("HEALTH_PORT")
	)

	log.Debugf("Starting server...")

	// main router
	router := apmhttprouter.New()
	router.NotFound = http.FileServer(
		http.Dir("static"),
	)

	cfg := app.New()
	cfg.Init()
	cfg.Routes(router)

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = middlewares.Logger(
		middlewares.I18(
			middlewares.Validate(router),
		),
	)

	errChan := make(chan error, 10)

	// health router
	healthServer := health.StartServer(healthAddr)

	go func() {
		log.Debugf("Health service listening on %s", healthAddr)
		errChan <- healthServer.ListenAndServe()
	}()
	go func() {
		log.Debugf("HTTP service listening on %s", httpAddr)
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				apm.CaptureError(context.TODO(), err).Send()
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Debugf(fmt.Sprintf("Captured %v. Exiting...", s))
			health.SetReadinessStatus(http.StatusServiceUnavailable)
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
