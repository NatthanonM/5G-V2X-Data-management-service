package controllers

import (
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
	"fmt"
	"time"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccidentController struct {
	*services.AccidentService
}

func NewAccidentController(accidentSrvc *services.AccidentService) *AccidentController {
	return &AccidentController{
		AccidentService: accidentSrvc,
	}
}

func (ac *AccidentController) CreateAccidentData(ctx context.Context, req *proto.AccidentData) (*proto.CreateAccidentDataResponse, error) {
	accidentID, err := ac.AccidentService.StoreData(
		req.Username,
		req.CarId,
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
	records, err := ac.AccidentService.GetRecords(req.From, req.To, req.CarId)

	var accidentList []*proto.AccidentData
	for _, elem := range records {
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
	return &proto.GetAccidentDataResponse{
		Accidents: accidentList,
	}, nil
}

// GetHourlyAccidentOfCurrentDay ...
func (ac *AccidentController) GetHourlyAccidentOfCurrentDay(ctx context.Context, req *proto.GetHourlyAccidentOfCurrentDayRequest) (*proto.GetHourlyAccidentOfCurrentDayResponse, error) {
	if req.Hour < 0 || req.Hour > 23 {
		err := status.Error(codes.InvalidArgument, "Hour must be between 0 to 23")
		fmt.Println(err)
		return nil, err
	}

	hourlyAccidentOfCurrentDay, err := ac.AccidentService.GetHourlyAccidentOfCurrentDay(req.Hour)

	var accidentList []*proto.AccidentData
	for _, elem := range hourlyAccidentOfCurrentDay {
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
	return &proto.GetHourlyAccidentOfCurrentDayResponse{
		Accidents: accidentList,
	}, nil
}

// GetGetNumberOfAccidentCurrentYearDetailDay
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
			Name:  	  elem.Name,
			Data:     elem.Data,
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

// GetGetNumberOfAccidentCurrentYearDetailDay
func (ac *AccidentController) GetNumberOfAccidentTimeBar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfAccidentTimeBarResponse, error) {
	year, month, day := time.Now().Date()

	numberOfAccidentTimeBar, err := ac.AccidentService.GetNumberOfAccidentTimeBar(day,int(month),year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetNumberOfAccidentTimeBarResponse{
		Accidents: numberOfAccidentTimeBar,
	}, nil
}

func (ac *AccidentController) GetNumberOfAccidentStreet(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfAccidentStreetResponse, error) {
	year, month, day := time.Now().Date()
	var no []int32
	var label []string
	acStreet, err := ac.AccidentService.GetNumberOfAccidentStreet(day,int(month),year,day,int(month),year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for k, v := range acStreet {
		label= append(label,string(k))
		no = append(no, int32(v))
	}
	anAccident := &proto.AccidentStatPieData{
		Series: no,
		Labels: label,		
	}
	return &proto.GetNumberOfAccidentStreetResponse{
		Accidents: anAccident,
	}, nil
}