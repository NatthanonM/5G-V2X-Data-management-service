package controllers

import (
	"5g-v2x-data-management-service/internal/services"
	"5g-v2x-data-management-service/internal/utils"
	proto "5g-v2x-data-management-service/pkg/api"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DrowsinessController struct {
	*services.DrowsinessService
}

func NewDrowsinessController(drowsinessSrvc *services.DrowsinessService) *DrowsinessController {
	return &DrowsinessController{
		DrowsinessService: drowsinessSrvc,
	}
}

func (dc *DrowsinessController) CreateDrowsinessData(ctx context.Context, req *proto.DrowsinessData) (*proto.CreateDrowsinessDataResponse, error) {
	drowsinessID, err := dc.DrowsinessService.StoreData(
		req.Username,
		req.CarId,
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
			Name:  	  elem.Name,
			Data:     elem.Data,
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
	numberOfDrowsinessTimeBar, err := ac.DrowsinessService.GetNumberOfDrowsinessTimeBar(day,int(month),year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetNumberOfDrowsinessTimeBarResponse{
		Drowsinesss: numberOfDrowsinessTimeBar,
	}, nil
}