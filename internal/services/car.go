package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CarService struct {
	CarRepository *repositories.CarRepository
}

func NewCarService(carRepository *repositories.CarRepository) *CarService {
	return &CarService{
		CarRepository: carRepository,
	}
}

func (cs *CarService) GetCarByVehicleRegistrationNumber(vehRegNo string) (*models.Car, error) {
	car, err := cs.CarRepository.FindOne(nil, &vehRegNo)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (cs *CarService) RegisterNewCar(car *models.Car) (*string, error) {
	if _, err := cs.GetCarByVehicleRegistrationNumber(*car.VehicleRegistrationNumber); err == nil {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Vehicle registration number `%s` is already existed.", *car.VehicleRegistrationNumber))
	}

	carID, err := cs.CarRepository.Create(car)
	if err != nil {
		return nil, err
	}
	return carID, nil
}

func (cs *CarService) GetAllCar() ([]*models.Car, error) {
	carList, err := cs.CarRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return carList, nil
}

func (cs *CarService) GetCar(carID string) (*models.Car, error) {
	car, err := cs.CarRepository.FindOne(&carID, nil)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Car not found.")
	}
	return car, nil
}

func (cs *CarService) UpdateCar(updateCar *models.Car) error {
	if updateCar.VehicleRegistrationNumber != nil {
		if car, err := cs.GetCarByVehicleRegistrationNumber(*updateCar.VehicleRegistrationNumber); err == nil && car.CarID != updateCar.CarID {
			return status.Error(codes.AlreadyExists, fmt.Sprintf("Vehicle registration number `%s` is already existed.", *updateCar.VehicleRegistrationNumber))
		}
	}

	_, err := cs.GetCar(updateCar.CarID)

	if err != nil {
		return status.Error(codes.NotFound, "Car not found.")
	}

	err = cs.CarRepository.Update(updateCar)
	if err != nil {
		return status.Error(codes.Internal, "Update car failed")
	}

	return err
}

func (cs *CarService) DeleteCar(carId string) error {
	_, err := cs.GetCar(carId)

	if err != nil {
		return status.Error(codes.NotFound, "Car not found.")
	}

	err = cs.CarRepository.Delete(carId)
	if err != nil {
		return status.Error(codes.Internal, "Delete car failed")
	}

	return err
}
