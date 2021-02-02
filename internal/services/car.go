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

func (cs *CarService) RegisterNewCar(car *models.Car) (*string, error) {
	carID, err := cs.CarRepository.Create(car)
	if err != nil {
		return nil, err
	}
	return carID, nil
}
