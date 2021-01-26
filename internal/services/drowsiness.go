package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"time"
)

type DrowsinessService struct {
	crud                 *repositories.CRUDRepository
	DrowsinessRepository *repositories.DrowsinessRepository
}

func NewDrowsinessService(crud *repositories.CRUDRepository, drowsinessRepository *repositories.DrowsinessRepository) *DrowsinessService {
	return &DrowsinessService{
		crud:                 crud,
		DrowsinessRepository: drowsinessRepository,
	}
}

func (ds *DrowsinessService) StoreData(username string, carID string, lat float64, lng float64, time time.Time, responseTime float64, workingHour float64) (string, error) {
	var drowsiness models.Drowsiness
	drowsiness.Username = username
	drowsiness.CarID = carID
	drowsiness.Latitude = lat
	drowsiness.Longitude = lng
	drowsiness.Time = time
	drowsiness.WorkingHour = workingHour
	drowsiness.ResponseTime = responseTime
	id, err := ds.crud.Create("drowsiness", &drowsiness)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (ds *DrowsinessService) GetHourlyDrowsinessOfCurrentDay(hour int32) ([]*models.Drowsiness, error) {
	result, err := ds.DrowsinessRepository.GetHourlyDrowsinessOfCurrentDay(hour)
	if err != nil {
		return nil, err
	}
	return result, nil
}
