package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"time"
)

type AccidentService struct {
	crud *repositories.CRUDRepository
}

func NewAccidentService(crud *repositories.CRUDRepository) *AccidentService {
	return &AccidentService{
		crud: crud,
	}
}

func (as *AccidentService) StoreData(username string) error {
	var accident models.Accident
	accident.Username = username
	accident.CarID = ""
	accident.Latitude = 0
	accident.Longitude = 0
	accident.Time = time.Unix(0, 0)

	_, err := as.crud.Create("accident", &accident)

	if err != nil {
		panic(err.Error())
	}
	return nil
}
