package models

import "time"

type Car struct {
	CarID                     string    `bson:"_id"`
	VehicleRegistrationNumber string    `bson:"vehicle_registration_number"`
	CarType                   string    `bson:"car_type"`
	RegisteredAt              time.Time `bson:"registered_at"`
}
