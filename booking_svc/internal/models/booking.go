package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"booking_id"`
	PickupLoc  Location           `json:"pickuploc" bson:"pickuploc"`
	Dropoff    Location           `json:"dropoff" bson:"dropoff"`
	Price      int                `json:"price" bson:"price"`
	RideStatus string             `json:"ride_status" bson:"ride_status"`
	DriverID   *string            `json:"driver_id,omitempty" bson:"driver_id,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

func NowUTC() time.Time {
	return time.Now().UTC()
}
