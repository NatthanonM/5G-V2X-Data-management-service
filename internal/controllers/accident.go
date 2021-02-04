package controllers

import (
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
	"fmt"

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

// GetAllAccidentData ...
func (ac *AccidentController) GetAccidentData(ctx context.Context, req *proto.GetAccidentDataRequest) (*proto.GetAccidentDataResponse, error) {
	records, err := ac.AccidentService.GetRecords(req.From.AsTime(), req.To.AsTime())

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
