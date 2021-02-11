package models

import "time"

type Accident struct {
	CarID     string    `bson:"car_id"`
	Username  string    `bson:"username"`
	Time      time.Time `bson:"time"`
	Latitude  float64   `bson:"latitude"`
	Longitude float64   `bson:"longitude"`
	Road      string    `bson:"road"`
}
type AccidentStatCal struct {
	Name string  `json:"name"`
	Data []int32 `json:"data"`
	Username  string
	CarID     string `bson:"car_id"`
	Time      time.Time
	Latitude  float64
	Longitude float64
}
type NumberOfAccidentRoad struct {
	ID    string `bson:"_id" json:"id"`
	Total int32 `json:"total"`
}
