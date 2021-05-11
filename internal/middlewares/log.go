package middlewares

import (
	"net/http"

	"github.com/fakovacic/aig/internal/log"
)

// Logger print requests in console
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s - [%s]\n", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	})
}
