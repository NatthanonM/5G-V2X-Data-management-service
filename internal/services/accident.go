package services

import (
	"5g-v2x-data-management-service/internal/models"
	"5g-v2x-data-management-service/internal/repositories"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (as *AccidentService) StoreData(username, carID, road string, lat float64, lng float64, time time.Time) (string, error) {
	var accident models.Accident
	accident.Username = username
	accident.CarID = carID
	accident.Latitude = lat
	accident.Longitude = lng
	accident.Time = time
	accident.Road = road

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

func (as *AccidentService) GetRecords(from, to *timestamppb.Timestamp, carID, username *string) ([]*models.Accident, error) {
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
				{"$gt", fromTime},
				{"$lte", toTime},
			},
		},
	}

	if carID != nil {
		filter = append(filter, bson.E{
			"car_id", carID,
		})
	}

	if username != nil {
		filter = append(filter, bson.E{
			"username", username,
		})
	}

	fmt.Println(fromTime, "-", toTime)
	fmt.Println(filter)

	result, err := as.AccidentRepository.Find(filter)

	for _, res := range result {
		fmt.Println(res.Time)
	}

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

func (as *AccidentService) GetNumberOfAccidentToCalendar(year int) ([]*models.AccidentStatCal, error) {
	result, err := as.AccidentRepository.GetNumberOfAccidentToCalendar(year)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *AccidentService) GetNumberOfAccidentTimeBar(day int, month int, year int) ([]int32, error) {
	result, err := as.AccidentRepository.GetNumberOfAccidentTimeBar(day, month, year)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (as *AccidentService) GetNumberOfAccidentStreet(startDay int, startMonth int, startYear int) (map[string]int32, error) {
	result, err := as.AccidentRepository.GetNumberOfAccidentStreet(startDay, startMonth, startYear)

	if err != nil {
		return nil, err
	}

	return result, nil
}
