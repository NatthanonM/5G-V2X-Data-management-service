package controllers

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/services"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if _, err := cc.CarService.GetCarByVehicleRegistrationNumber(req.VehicleRegistrationNumber); err == nil {
		return nil, status.Error(codes.AlreadyExists, "Vehicle registration number is already existed.")
	}
	car := models.Car{
		CarDetail:                 req.CarDetail,
		VehicleRegistrationNumber: req.VehicleRegistrationNumber,
	}
	carID, err := cc.CarService.RegisterNewCar(&car)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return &proto.RegisterNewCarResponse{
		CarId: *carID,
	}, nil
}
