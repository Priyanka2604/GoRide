package models

type Driver struct {
	ID        string `json:"driver_id" bson:"driver_id"`
	Name      string `json:"name" bson:"name"`
	Available bool   `json:"available" bson:"available"`
}
