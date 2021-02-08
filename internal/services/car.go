package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"

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
	filter := make(map[string]interface{})
	filter["vehicle_registration_number"] = vehRegNo
	car, err := cs.CarRepository.FindOne(filter)
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (cs *CarService) RegisterNewCar(car *models.Car) (*string, error) {
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
	filter := make(map[string]interface{})
	filter["_id"] = carID
	car, err := cs.CarRepository.FindOne(filter)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Car not found.")
	}
	return car, nil
}
