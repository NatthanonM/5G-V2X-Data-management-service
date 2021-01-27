package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"time"
)

type AccidentService struct {
	crud               *repositories.CRUDRepository
	AccidentRepository *repositories.AccidentRepository
}

func NewAccidentService(crud *repositories.CRUDRepository, accidentRepository *repositories.AccidentRepository) *AccidentService {
	return &AccidentService{
		crud:               crud,
		AccidentRepository: accidentRepository,
	}
}

func (as *AccidentService) StoreData(username string, carID string, lat float64, lng float64, time time.Time) (string, error) {
	var accident models.Accident
	accident.Username = username
	accident.CarID = carID
	accident.Latitude = lat
	accident.Longitude = lng
	accident.Time = time

	id, err := as.crud.Create("accident", &accident)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (as *AccidentService) GetAllRecords() ([]*models.Accident, error) {

	result, err := as.AccidentRepository.FindAll()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (as *AccidentService) GetHourlyAccidentOfCurrentDay(hour int32) ([]*models.Accident, error) {

	result, err := as.AccidentRepository.GetHourlyAccidentOfCurrentDay(hour)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (as *AccidentService) GetNumberOfAccidentCurrentYear() ([]*models.AccidentStatCal, error) {
	result, err := as.AccidentRepository.GetNumberOfAccidentCurrentYear()

	if err != nil {
		return nil, err
	}

	return result, nil
}