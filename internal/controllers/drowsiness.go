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

func (dc *DrowsinessController) GetDrowsinessData(ctx context.Context, req *proto.GetDrowsinessDataRequest) (*proto.GetDrowsinessDataResponse, error) {
	drowsinessData, err := dc.DrowsinessService.GetDrowsiness(req.From, req.To, req.CarId, req.Username)
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
			Road:         drowsiness.Road,
		})
	}
	return &proto.GetDrowsinessDataResponse{
		Drowsinesses: drowsinesses,
	}, nil
}

func (dc *DrowsinessController) GetNumberOfDrowsinessToCalendar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfDrowsinessToCalendarResponse, error) {
	year := time.Now().Year()
	numberOfDrowsinessCurrentYear, err := dc.DrowsinessService.GetNumberOfDrowsinessToCalendar(year)
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

func (dc *DrowsinessController) GetNumberOfDrowsinessTimeBar(ctx context.Context, req *empty.Empty) (*proto.GetNumberOfDrowsinessTimeBarResponse, error) {
	year, month, day := time.Now().Date()
	numberOfDrowsinessTimeBar, err := dc.DrowsinessService.GetNumberOfDrowsinessTimeBar(day, int(month), year)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.GetNumberOfDrowsinessTimeBarResponse{
		Drowsinesss: numberOfDrowsinessTimeBar,
	}, nil
}

func (dc *DrowsinessController) GetDrowsinessStatGroupByHour(ctx context.Context, req *proto.GetDrowsinessStatGroupByHourRequest) (*proto.GetDrowsinessStatGroupByHourResponse, error) {
	countEachHour, err := dc.DrowsinessService.GetDrowsinessStatGroupByHour(req.From, req.To, req.DriverUsername)
	if err != nil {
		return nil, err
	}
	drowsinesses := []int64{}
	for _, e := range countEachHour {
		drowsinesses = append(drowsinesses, int64(e))
	}
	return &proto.GetDrowsinessStatGroupByHourResponse{
		Drowsinesses: drowsinesses,
	}, nil
}
