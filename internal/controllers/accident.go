package controllers

import (
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

type AccidentController struct {
	*services.AccidentService
	*services.GoogleService
}

func NewAccidentController(accidentSrvc *services.AccidentService, googleSrvc *services.GoogleService) *AccidentController {
	return &AccidentController{
		AccidentService: accidentSrvc,
		GoogleService:   googleSrvc,
	}
}

func (ac *AccidentController) CreateAccidentData(ctx context.Context, req *proto.AccidentData) (*proto.CreateAccidentDataResponse, error) {
	fmt.Println(req)
	road, err := ac.GoogleService.ReverseGeocoding(req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}
	fmt.Println(*road)

	accidentID, err := ac.AccidentService.StoreData(
		req.Username,
		req.CarId,
		*road,
		req.Latitude,
		req.Longitude,
		*utils.WrapperTimeStamp(req.Time),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.CreateAccidentDataResponse{
		AccidentId: accidentID,
	}, nil
}

// GetAllAccidentData ...
func (ac *AccidentController) GetAllAccidentData(ctx context.Context, req *empty.Empty) (*proto.GetAllAccidentDataResponse, error) {
	allRecords, err := ac.AccidentService.GetAllRecords()

	var accidentList []*proto.AccidentData
	for _, elem := range allRecords {
		anAccident := proto.AccidentData{
			Username:  elem.Username,
			CarId:     elem.CarID,
			Latitude:  elem.Latitude,
			Longitude: elem.Longitude,
			Time:      utils.WrapperTime(&elem.Time),
		}
		accidentList = append(accidentList, &anAccident)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetAllAccidentDataResponse{
		Accidents: accidentList,
	}, nil
}

// GetAccidentData ...
func (ac *AccidentController) GetAccidentData(ctx context.Context, req *proto.GetAccidentDataRequest) (*proto.GetAccidentDataResponse, error) {
	records, err := ac.AccidentService.GetRecords(req.From, req.To, req.CarId, req.Username)

	var accidentList []*proto.AccidentData
	for _, elem := range records {
		anAccident := proto.AccidentData{
			Username:  elem.Username,
			CarId:     elem.CarID,
			Latitude:  elem.Latitude,
			Longitude: elem.Longitude,
			Time:      utils.WrapperTime(&elem.Time),
			Road:      elem.Road,
		}
		accidentList = append(accidentList, &anAccident)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetAccidentDataResponse{
		Accidents: accidentList,
	}, nil
}

// GetHourlyAccidentOfCurrentDay ...
// func (ac *AccidentController) GetHourlyAccidentOfCurrentDay(ctx context.Context, req *proto.GetHourlyAccidentOfCurrentDayRequest) (*proto.GetHourlyAccidentOfCurrentDayResponse, error) {
// 	if req.Hour < 0 || req.Hour > 23 {
// 		err := status.Error(codes.InvalidArgument, "Hour must be between 0 to 23")
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	hourlyAccidentOfCurrentDay, err := ac.AccidentService.GetHourlyAccidentOfCurrentDay(req.Hour)

// 	var accidentList []*proto.AccidentData
// 	for _, elem := range hourlyAccidentOfCurrentDay {
// 		anAccident := proto.AccidentData{
// 			Username:  elem.Username,
// 			CarId:     elem.CarID,
// 			Latitude:  elem.Latitude,
// 			Longitude: elem.Longitude,
// 			Time:      utils.WrapperTime(&elem.Time),
// 			Road:      elem.Road,
// 		}
// 		accidentList = append(accidentList, &anAccident)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return &proto.GetHourlyAccidentOfCurrentDayResponse{
// 		Accidents: accidentList,
// 	}, nil
// }

// GetNumberOfAccidentToCalendar
func (ac *AccidentController) GetNumberOfAccidentToCalendar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfAccidentToCalendarResponse, error) {
	year := time.Now().Year()
	numberOfAccidentCurrentYear, err := ac.AccidentService.GetNumberOfAccidentToCalendar(year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var accidentList []*proto.AccidentStatCalData
	for _, elem := range numberOfAccidentCurrentYear {
		anAccident := proto.AccidentStatCalData{
			Name: elem.Name,
			Data: elem.Data,
		}
		accidentList = append(accidentList, &anAccident)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.GetNumberOfAccidentToCalendarResponse{
		Accidents: accidentList,
	}, nil
}

// GetNumberOfAccidentTimeBar
func (ac *AccidentController) GetNumberOfAccidentTimeBar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfAccidentTimeBarResponse, error) {
	// year, month, day := time.Now().Date()

	numberOfAccidentTimeBar, err := ac.AccidentService.GetNumberOfAccidentTimeBar(1, 0, 1970)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetNumberOfAccidentTimeBarResponse{
		Accidents: numberOfAccidentTimeBar,
	}, nil
}

func (ac *AccidentController) GetNumberOfAccidentStreet(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfAccidentStreetResponse, error) {
	year, month, day := time.Now().UTC().Date()
	var no []int32
	var label []string
	acStreet, err := ac.AccidentService.GetNumberOfAccidentStreet(day, int(month), year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, elem := range acStreet {
		if elem.ID == "" {
			label = append(label, "N/A")
		} else {
			label = append(label, elem.ID)
		}
		no = append(no, elem.Total)
	}
	anAccident := &proto.AccidentStatPieData{
		Series: no,
		Labels: label,
	}
	return &proto.GetNumberOfAccidentStreetResponse{
		Accidents: anAccident,
	}, nil
}

func (ac *AccidentController) GetAccidentStatGroupByHour(ctx context.Context, req *proto.GetAccidentStatGroupByHourRequest) (*proto.GetAccidentStatGroupByHourResponse, error) {
	countEachHour, err := ac.AccidentService.GetAccidentStatGroupByHour(req.From, req.To, req.DriverUsername)
	if err != nil {
		return nil, err
	}
	accidents := []int64{}
	for _, e := range countEachHour {
		accidents = append(accidents, int64(e))
	}
	return &proto.GetAccidentStatGroupByHourResponse{
		Accidents: accidents,
	}, nil
}

func (ac *AccidentController) GetTopNRoad(ctx context.Context, req *proto.GetTopNRoadRequest) (*proto.GetTopNRoadResponse, error) {
	topNRoadResult, err := ac.AccidentService.FindTopNRoad(req.From, req.To, req.N)
	// fmt.Println(req.From)
	// fmt.Println(req.To)
	// fmt.Println(*req.N)
	if err != nil {
		return nil, err
	}
	var topNRoad []*proto.TopNRoad
	for _, v := range topNRoadResult {
		topNRoad = append(topNRoad, &proto.TopNRoad{
			RoadName:      v.RoadName,
			AccidentCount: v.AccidentCount,
		})
	}
	return &proto.GetTopNRoadResponse{
		TopNRoad: topNRoad,
	}, nil
}
