package endpoints

import "github.com/fakovacic/aig/internal/types/endpoints"

const (
	// Register const
	Register endpoints.Endpoint = "/register"

	// Dashoard const
	Dashoard endpoints.Endpoint = "/"

	// ControllerCommand const
	ControllerCommand endpoints.Endpoint = "/controller/:id/:command/:power"
)
