package models

import (
	"gorm.io/gorm"
)

type Locations struct {
	Address    string   `json:"address"`
	CityID     uint     `json:"city_id" gorm:"type:uuid"`
	City       City     `json:"city" gorm:"foreignKey:CityID"`
	ProvinceID uint     `json:"province_id" gorm:"type:uuid"`
	Province   Province `json:"province" gorm:"foreignKey:ProvinceID"`
	PostalCode string   `json:"postal_code"`
	Country    string   `json:"country"`
	Geo        string   `json:"geo"`
}

type City struct {
	gorm.Model
	ID       uint   `json:"ID" gorm:"type:uuid;primaryKey"`
	CityName string `json:"city_name"`
}

type Province struct {
	gorm.Model
	ID           uint   `json:"ID" gorm:"type:uuid;primaryKey"`
	ProvinceName string `json:"province_name"`
}
