package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyCode string `json:"code" gorm:"primaryKey"`
	CompanyName string `json:"name"`
}