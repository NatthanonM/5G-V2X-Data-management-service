package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
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
	carList, err := cs.CarRepository.FineAll()
	if err != nil {
		return nil, err
	}
	return carList, nil
}
