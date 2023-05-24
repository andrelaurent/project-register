package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Locations struct {
	gorm.Model
	ID         uint     `json:"ID"`
	Address    string   `json:"address"`
	City       City     `json:"city"`
	Province   Province `json:"province"`
	PostalCode string   `json:"postal_code"`
	Country    string   `json:"country"`
	Geo        string   `json:"geo"`
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