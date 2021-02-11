package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (ds *DrowsinessService) StoreData(username, carID, road string, lat float64, lng float64, time time.Time, responseTime float64, workingHour float64) (string, error) {

	var drowsiness models.Drowsiness
	drowsiness.Username = username
	drowsiness.CarID = carID
	drowsiness.Latitude = lat
	drowsiness.Longitude = lng
	drowsiness.Time = time
	drowsiness.WorkingHour = workingHour
	drowsiness.ResponseTime = responseTime
	drowsiness.Road = road
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

func (ds *DrowsinessService) GetDrowsiness(carID, username *string) ([]*models.Drowsiness, error) {
	filter := bson.D{{}}
	if carID != nil {
		filter = append(filter, bson.E{
			"car_id", *carID,
		})
	}
	if username != nil {
		filter = append(filter, bson.E{
			"username", *username,
		})
	}
	drowsinessData, err := ds.DrowsinessRepository.Find(filter)
	if err != nil {
		return nil, err
	}
	return drowsinessData, nil
}

func (as *DrowsinessService) GetNumberOfDrowsinessToCalendar(year int) ([]*models.DrowsinessStatCal, error) {
	result, err := as.DrowsinessRepository.GetNumberOfDrowsinessToCalendar(year)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *DrowsinessService) GetNumberOfDrowsinessTimeBar(day int, month int, year int) ([]int32, error) {
	result, err := as.DrowsinessRepository.GetNumberOfDrowsinessTimeBar(day, month, year)
	if err != nil {
		return nil, err
	}

	return result, nil
}
