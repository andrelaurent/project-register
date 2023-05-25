package models

import (
	"gorm.io/gorm"
)

type Locations struct {
	gorm.Model
	ClientID   uint     `json:"client_id"`
	Address    string   `json:"address"`
	CityID     uint     `json:"city_id"`
	City       City     `json:"city"`
	ProvinceID uint     `json:"province_id"`
	Province   Province `json:"province"`
	PostalCode string   `json:"postal_code"`
	Country    string   `json:"country"`
	Geo        string   `json:"geo"`
}

type City struct {
	gorm.Model
	ID       uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	CityName string `json:"city_name"`
}

type Province struct {
	gorm.Model
	ID           uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProvinceName string `json:"province_name"`
}
