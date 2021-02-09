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
