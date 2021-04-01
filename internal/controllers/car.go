package controllers

import (
	"5g-v2x-data-management-service/internal/config"
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
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
		MfgAt:                     req.MfgAt.AsTime(),
	}
	carID, err := cc.CarService.RegisterNewCar(&car)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return &proto.RegisterNewCarResponse{
		CarId: *carID,
	}, nil
}

func (cc *CarController) GetCarList(ctx context.Context, req *empty.Empty) (*proto.GetCarListResponse, error) {
	carList, err := cc.CarService.GetAllCar()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	var grpcCarList []*proto.Car
	for _, car := range carList {
		grpcCarList = append(grpcCarList, &proto.Car{
			CarId:                     car.CarID,
			VehicleRegistrationNumber: car.VehicleRegistrationNumber,
			CarDetail:                 car.CarDetail,
			RegisteredAt:              utils.WrapperTime(&car.RegisteredAt),
			MfgAt:                     utils.WrapperTime(&car.MfgAt),
		})
	}
	return &proto.GetCarListResponse{
		CarList: grpcCarList,
	}, nil
}

func (cc *CarController) GetCar(ctx context.Context, req *proto.GetCarRequest) (*proto.Car, error) {
	car, err := cc.CarService.GetCar(req.CarId)
	if err != nil {
		return nil, err
	}
	return &proto.Car{
		CarId:                     car.CarID,
		VehicleRegistrationNumber: car.VehicleRegistrationNumber,
		CarDetail:                 car.CarDetail,
		RegisteredAt:              utils.WrapperTime(&car.RegisteredAt),
		MfgAt:                     utils.WrapperTime(&car.MfgAt),
	}, nil
}

func (cc *CarController) UpdateCar(ctx context.Context, req *proto.UpdateCarRequest) (*proto.UpdateCarResponse, error) {
	err := cc.CarService.UpdateCar(&models.Car{
		CarID:                     req.CarId,
		CarDetail:                 req.CarDetail,
		VehicleRegistrationNumber: req.VehicleRegistrationNumber,
	})
	if err != nil {
		return nil, err
	}
	return &proto.UpdateCarResponse{}, nil
}

func (cc *CarController) DeleteCar(ctx context.Context, req *proto.DeleteCarRequest) (*proto.DeleteCarResponse, error) {
	// err := cc.CarService.UpdateCar(&models.Car{
	// 	CarID:                     req.CarId,
	// 	CarDetail:                 req.CarDetail,
	// 	VehicleRegistrationNumber: req.VehicleRegistrationNumber,
	// })
	// if err != nil {
	// 	return nil, err
	// }
	return &proto.DeleteCarResponse{}, nil
}
