package models

import "time"

type Drowsiness struct {
	CarID        string    `bson:"car_id"`
	Username     string    `bson:"username"`
	Time         time.Time `bson:"time"`
	ResponseTime float64   `bson:"response_time"`
	WorkingHour  float64   `bson:"working_hour"`
	Latitude     float64   `bson:"latitude"`
	Longitude    float64   `bson:"longitude"`
	Road         string    `bson:"road"`
}

type DrowsinessStatCal struct {
	Name string  `json:"name"`
	Data []int32 `json:"data"`
}
type NumberOfDrowsiness struct {
	ID    NumberOfDrowsinessHourField `bson:"_id" json:"id"`
	Total int32                       `json:"total"`
}
type NumberOfDrowsinessHourField struct {
	Hour int32 `json:"h"`
}
type NumberOfDrowsinessDate struct {
	ID    NumberOfDrowsinessDateField `bson:"_id" json:"id"`
	Total int32                     `json:"total"`
}
type NumberOfDrowsinessDateField struct {
    Date  string `bson:"date" json:"date"`
}
