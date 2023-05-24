package models

import (
	"github.com/andrelaurent/project-register/vendor/github.com/google/uuid"
	"gorm.io/gorm"
)

type Locations struct {
	gorm.Model
	Address    string    `json:"address"`
	City       uuid.UUID `json:"city"`
	Province   uuid.UUID `json:"province"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	Geo        string    `json:"geo"`
}

type City struct {
	gorm.Model
	ID       uuid.UUID `json:"ID" gorm:"primaryKey;"`
	CityName string    `json:"city_name"`
}

type Province struct {
	gorm.Model
	ID           uuid.UUID `json:"ID" gorm:"primaryKey;"`
	ProvinceName string    `json:"province_name"`
}
