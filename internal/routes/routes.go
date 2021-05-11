package routes

import (
	"github.com/fakovacic/aig/internal/config"
	"github.com/fakovacic/aig/internal/types/methods"
	"github.com/julienschmidt/httprouter"
	"go.elastic.co/apm/module/apmhttprouter"
)

type Endpoint interface {
	String() string
}

type Method interface {
	String() string
}

type RoutesMap []map[Endpoint]map[Method]Route

// Register register routes
func Register(router *apmhttprouter.Router, rm RoutesMap) {
	for i := range rm {
		for route, methods := range rm[i] {
			for method, handler := range methods {
				handler.Register(router, route, method)
			}
		}
	}
}

// Route route struct
type Route struct {
	Handler httprouter.Handle
	Config  config.Config
}

// Register route
func (r *Route) Register(router *apmhttprouter.Router, endpoint Endpoint, method Method) {
	switch method {
	case methods.Get:
		router.GET(endpoint.String(), r.H())
	case methods.Post:
		router.POST(endpoint.String(), r.H())
	}
}

// H check handler
func (r *Route) H() httprouter.Handle {
	return r.Handler
}
