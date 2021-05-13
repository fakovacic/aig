package endpoints

import "github.com/fakovacic/aig/internal/types/endpoints"

const (
	// Register const
	Register endpoints.Endpoint = "/register"

	// Dashoard const
	Dashoard endpoints.Endpoint = "/"

	// Controller const
	Controller endpoints.Endpoint = "/controller/:id"

	// ControllerConfigAdd const
	ControllerConfigAdd endpoints.Endpoint = "/controller/:id/config/add"

	// ControllerConfigRemove const
	ControllerConfigRemove endpoints.Endpoint = "/controller/:id/config/remove/:config"
)
