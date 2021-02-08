package models

import "time"

type Accident struct {
	CarID     string    `bson:"car_id"`
	Username  string    `bson:"username"`
	Time      time.Time `bson:"time"`
	Latitude  float64   `bson:"latitude"`
	Longitude float64   `bson:"longitude"`
}
