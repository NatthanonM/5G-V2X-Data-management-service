package models

import "time"

type Drowsiness struct {
	Username     string
	CarID        string
	Time         time.Time
	ResponseTime float64
	WorkingHour  float64
	Latitude     float64
	Longitude    float64
}

type DrowsinessStatCal struct {
	Name     string     `json:"name"`
	Data	 []int32	`json:"data"`
}
