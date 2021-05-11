package middlewares

import (
	"net/http"

	"github.com/fakovacic/aig/internal/log"
	"github.com/fakovacic/aig/pkg/i18"
	"go.elastic.co/apm"
)

// I18 set lang
func I18(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := apm.TransactionFromContext(r.Context())
		s := t.StartSpan("en", "i18", nil)
		err := i18.ParseJSON("en")
		if err != nil {
			log.Errorf("parse i18:%v", err)
		}
		s.End()
		h.ServeHTTP(w, r)
	})
}
