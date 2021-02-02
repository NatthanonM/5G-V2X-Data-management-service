package controllers

import "go.uber.org/dig"

type Controller struct {
	*ControllerGateway
}

type ControllerGateway struct {
	dig.In
	*AccidentController
	*DrowsinessController
	*CarController
}

func NewController(cg *ControllerGateway) *Controller {
	return &Controller{
		ControllerGateway: cg,
	}
}
