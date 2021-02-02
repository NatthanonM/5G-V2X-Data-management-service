package controllers

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/services"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
)

type CarController struct {
	*services.CarService
	*config.Config
}

func NewCarController(CarService *services.CarService, Config *config.Config) *CarController {
	return &CarController{
		CarService: CarService,
		Config:     Config,
	}
}

func (cc *CarController) RegisterNewCar(ctx context.Context, req *proto.RegisterNewCarRequest) (*proto.RegisterNewCarResponse, error) {
	car := models.Car{
		CarType:                   req.CarType,
		VehicleRegistrationNumber: req.VehicleRegistrationNumber,
	}
	carID, err := cc.CarService.RegisterNewCar(&car)
	if err != nil {
		return nil, err
	}
	return &proto.RegisterNewCarResponse{
		CarId: *carID,
	}, nil
}
