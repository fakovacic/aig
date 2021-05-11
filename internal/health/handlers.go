package health

import (
	"net/http"
	"sync"
)

var (
	healthzStatus   = http.StatusOK
	readinessStatus = http.StatusOK
	mu              sync.RWMutex
)

func GetHealthStatus() int {
	mu.RLock()
	defer mu.RUnlock()
	return healthzStatus
}

func GetReadinessStatus() int {
	mu.RLock()
	defer mu.RUnlock()
	return readinessStatus
}

func SetHealthStatus(status int) {
	mu.Lock()
	healthzStatus = status
	mu.Unlock()
}

func SetReadinessStatus(status int) {
	mu.Lock()
	readinessStatus = status
	mu.Unlock()
}

// HealthzHandler responds to health check requests.
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(GetHealthStatus())
}

// ReadinessHandler responds to readiness check requests.
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(GetReadinessStatus())
}

func ReadinessStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch GetReadinessStatus() {
	case http.StatusOK:
		SetReadinessStatus(http.StatusServiceUnavailable)
	case http.StatusServiceUnavailable:
		SetReadinessStatus(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
}

func HealthStatusHandler(w http.ResponseWriter, r *http.Request) {
	switch GetHealthStatus() {
	case http.StatusOK:
		SetHealthStatus(http.StatusServiceUnavailable)
	case http.StatusServiceUnavailable:
		SetHealthStatus(http.StatusOK)
	}
	w.WriteHeader(http.StatusOK)
}
