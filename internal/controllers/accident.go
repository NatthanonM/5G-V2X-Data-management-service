package controllers

import (
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
)

type AccidentController struct {
	*services.AccidentService
}

func NewAccidentController(accidentSrvc *services.AccidentService) *AccidentController {
	return &AccidentController{
		AccidentService: accidentSrvc,
	}
}

func (ac *AccidentController) CreateAccidentData(ctx context.Context, req *proto.AccidentData) (*proto.CreateAccidentDataResponse, error) {
	accidentID, err := ac.AccidentService.StoreData(
		req.Username,
		req.CarId,
		req.Latitude,
		req.Longitude,
		*utils.WrapperTimeStamp(req.Time),
	)
	if err != nil {
		return nil, err
	}
	return &proto.CreateAccidentDataResponse{
		AccidentId: accidentID,
	}, nil
}
