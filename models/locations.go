package models

import (
	"gorm.io/gorm"
)

type Locations struct {
	gorm.Model
	ID         uint     `json:"ID" gorm:"primaryKey:autoIncrement"`
	Address    string   `json:"address"`
	CityID     uint     `json:"city_id" gorm:"index"`
	City       City     `json:"city"`
	ProvinceID uint     `json:"province_id" gorm:"index"`
	Province   Province `json:"province"`
	PostalCode string   `json:"postal_code"`
	Country    string   `json:"country"`
	Geo        string   `json:"geo"`
	ClientID   *uint    `json:"client_id"`
	ContactID  *uint    `json:"contact_id"`
}

type City struct {
	gorm.Model
	ID         uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	CityName   string `json:"city_name"`
	ProvinceID *uint  `json:"province_id"`
}

type Province struct {
	gorm.Model
	ID           uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProvinceName string `json:"province_name"`
	Cities       []City `json:"cities"`
}
