package models

import "time"

type Accident struct {
	Username  string
	CarID     string
	Time      time.Time
	Latitude  float64
	Longitude float64
}
