package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (ds *DrowsinessService) GetDrowsiness(from, to *timestamppb.Timestamp, carID, username *string) ([]*models.Drowsiness, error) {
	fromTime := time.Date(1970, time.Month(0), 0, 0, 0, 0, 0, time.UTC)
	toTime := time.Now()
	if from != nil {
		fromTime = from.AsTime()
	}
	if to != nil {
		toTime = to.AsTime()
	}

	filter := bson.D{
		{
			"time", bson.D{
				{"$gte", fromTime},
				{"$lt", toTime},
			},
		},
	}

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

func (ds *DrowsinessService) GetNumberOfDrowsinessToCalendar(year int64) ([]*models.DrowsinessStatCal, error) {
	result, err := ds.DrowsinessRepository.GetNumberOfDrowsinessToCalendar(year)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds *DrowsinessService) GetNumberOfDrowsinessTimeBar(from, to time.Time) ([]int32, error) {
	result, err := ds.DrowsinessRepository.GetNumberOfDrowsinessTimeBar(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ds *DrowsinessService) GetDrowsinessStatGroupByHour(from, to *timestamppb.Timestamp, driverUsername *string) ([24]int32, error) {
	result, err := ds.DrowsinessRepository.GetDrowsinessStatGroupByHour(from, to, driverUsername)
	if err != nil {
		return [24]int32{}, err
	}
	return result, nil
}
