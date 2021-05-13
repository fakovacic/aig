package routes

import (
	"fmt"

	"github.com/fakovacic/aig/internal/app/actions"
	"github.com/fakovacic/aig/internal/app/config"
	"github.com/fakovacic/aig/internal/app/endpoints"
	controllers "github.com/fakovacic/aig/internal/app/handlers"
	"github.com/fakovacic/aig/internal/handlers"
	"github.com/fakovacic/aig/internal/routes"
	"github.com/fakovacic/aig/internal/types/methods"
	"github.com/fakovacic/aig/pkg/templates"
)

// Routes routes for sources service
func Routes(cfg config.Config) map[routes.Endpoint]map[routes.Method]routes.Route {
	base := templates.New(
		"index.html",
		nil,
		"templates/index.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/meta.html",
		"templates/notifications.html",
	)
	err := base.Parse()
	if err != nil {
		fmt.Printf("parsing controllers base error:%v", err)
	}

	h := controllers.New(cfg)

	return map[routes.Endpoint]map[routes.Method]routes.Route{
		endpoints.Dashoard: {
			methods.Get: routes.Route{
				Config: cfg,
				Handler: handlers.Handle(
					templates.New(
						"dashboard.html",
						base.Template,
						"templates/app/dashboard.html",
					), h, actions.Dashoard,
				),
			},
		},
		endpoints.Register: {
			methods.Post: routes.Route{
				Config:  cfg,
				Handler: handlers.Handle(nil, h, actions.Register),
			},
		},
		endpoints.Controller: {
			methods.Post: routes.Route{
				Config:  cfg,
				Handler: handlers.Handle(nil, h, actions.Controller),
			},
		},
		endpoints.ControllerConfigAdd: {
			methods.Get: routes.Route{
				Config:  cfg,
				Handler: handlers.Handle(nil, h, actions.ControllerConfigAdd),
			},
		},
		endpoints.ControllerConfigRemove: {
			methods.Get: routes.Route{
				Config:  cfg,
				Handler: handlers.Handle(nil, h, actions.ControllerConfigRemove),
			},
		},
	}
}
