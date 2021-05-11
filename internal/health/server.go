package health

import (
	"net/http"

	"github.com/braintree/manners"
	"github.com/fakovacic/aig/internal/middlewares"
)

// StartServer starts health server
func StartServer(healthAddr string) *manners.GracefulServer {
	// health router
	hmux := http.NewServeMux()
	hmux.HandleFunc("/healthz", HealthzHandler)
	hmux.HandleFunc("/readiness", ReadinessHandler)
	hmux.HandleFunc("/healthz/status", HealthStatusHandler)
	hmux.HandleFunc("/readiness/status", ReadinessStatusHandler)
	healthServer := manners.NewServer()
	healthServer.Addr = healthAddr
	healthServer.Handler = middlewares.Logger(hmux)

	return healthServer
}
