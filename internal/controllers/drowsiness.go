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

type DrowsinessController struct {
	*services.DrowsinessService
	*services.GoogleService
}

func NewDrowsinessController(drowsinessSrvc *services.DrowsinessService, googleSrvc *services.GoogleService) *DrowsinessController {
	return &DrowsinessController{
		DrowsinessService: drowsinessSrvc,
		GoogleService:     googleSrvc,
	}
}

func (dc *DrowsinessController) CreateDrowsinessData(ctx context.Context, req *proto.DrowsinessData) (*proto.CreateDrowsinessDataResponse, error) {
	road, err := dc.GoogleService.ReverseGeocoding(req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}

	drowsinessID, err := dc.DrowsinessService.StoreData(
		req.Username,
		req.CarId,
		*road,
		req.Latitude,
		req.Longitude,
		*utils.WrapperTimeStamp(req.Time),
		req.ResponseTime,
		req.WorkingHour,
	)
	if err != nil {
		return nil, err
	}
	return &proto.CreateDrowsinessDataResponse{
		DrowsinessId: drowsinessID,
	}, nil
}

func (dc *DrowsinessController) GetHourlyDrowsinessOfCurrentDay(ctx context.Context, req *proto.GetHourlyDrowsinessOfCurrentDayRequest) (*proto.GetHourlyDrowsinessOfCurrentDayResponse, error) {
	if req.Hour < 0 || req.Hour > 23 {
		err := status.Error(codes.InvalidArgument, "Hour must be between 0 to 23")
		fmt.Println(err)
		return nil, err
	}

	hourlyDrowsinessOfCurrentDay, err := dc.DrowsinessService.GetHourlyDrowsinessOfCurrentDay(req.Hour)

	var drosinessList []*proto.DrowsinessData
	for _, elem := range hourlyDrowsinessOfCurrentDay {
		drosinessList = append(drosinessList, &proto.DrowsinessData{
			Username:     elem.Username,
			CarId:        elem.CarID,
			Latitude:     elem.Latitude,
			Longitude:    elem.Longitude,
			Time:         utils.WrapperTime(&elem.Time),
			ResponseTime: elem.ResponseTime,
			WorkingHour:  elem.WorkingHour,
		})
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetHourlyDrowsinessOfCurrentDayResponse{
		Drowsinesses: drosinessList,
	}, nil
}

func (dc *DrowsinessController) GetDrowsinessData(ctx context.Context, req *proto.GetDrowsinessDataRequest) (*proto.GetDrowsinessDataResponse, error) {
	if req.CarId == nil && req.Username == nil {
		err := status.Error(codes.InvalidArgument, "Car id or username must be provided")
		return nil, err
	}
	drowsinessData, err := dc.DrowsinessService.GetDrowsiness(req.CarId, req.Username)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	drowsinesses := []*proto.DrowsinessData{}
	for _, drowsiness := range drowsinessData {
		drowsinesses = append(drowsinesses, &proto.DrowsinessData{
			CarId:        drowsiness.CarID,
			Username:     drowsiness.Username,
			Time:         utils.WrapperTime(&drowsiness.Time),
			ResponseTime: drowsiness.ResponseTime,
			WorkingHour:  drowsiness.WorkingHour,
			Latitude:     drowsiness.Latitude,
			Longitude:    drowsiness.Longitude,
		})
	}
	return &proto.GetDrowsinessDataResponse{
		Drowsinesses: drowsinesses,
	}, nil
}

func (ac *DrowsinessController) GetNumberOfDrowsinessToCalendar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfDrowsinessToCalendarResponse, error) {
	year := time.Now().Year()
	numberOfDrowsinessCurrentYear, err := ac.DrowsinessService.GetNumberOfDrowsinessToCalendar(year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var drowsinessList []*proto.DrowsinessStatCalData
	for _, elem := range numberOfDrowsinessCurrentYear {
		anDrowsiness := proto.DrowsinessStatCalData{
			Name: elem.Name,
			Data: elem.Data,
		}
		drowsinessList = append(drowsinessList, &anDrowsiness)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &proto.GetNumberOfDrowsinessToCalendarResponse{
		Drowsinesss: drowsinessList,
	}, nil
}

func (ac *DrowsinessController) GetNumberOfDrowsinessTimeBar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfDrowsinessTimeBarResponse, error) {
	year, month, day := time.Now().Date()
	numberOfDrowsinessTimeBar, err := ac.DrowsinessService.GetNumberOfDrowsinessTimeBar(day, int(month), year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetNumberOfDrowsinessTimeBarResponse{
		Drowsinesss: numberOfDrowsinessTimeBar,
	}, nil
}
