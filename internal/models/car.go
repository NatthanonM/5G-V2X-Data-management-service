package models

import "time"

type Car struct {
	CarID                     string     `bson:"_id"`
	VehicleRegistrationNumber *string    `bson:"vehicle_registration_number"`
	CarDetail                 *string    `bson:"car_detail"`
	RegisteredAt              time.Time  `bson:"registered_at"`
	MfgAt                     time.Time  `bson:"mfg_at"`
	DeletedAt                 *time.Time `bson:"deleted_at"`
}
