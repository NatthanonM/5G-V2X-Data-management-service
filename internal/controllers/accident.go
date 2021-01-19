package controllers

import "5g-v2x-data-management-service/internal/services"

type AccidentController struct {
	*services.AccidentService
}

func NewAccidentController(accidentSrvc *services.AccidentService) *AccidentController {
	return &AccidentController{
		AccidentService: accidentSrvc,
	}
}
