package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyID string `json:"code" gorm:"primaryKey"`
	CompanyName string `json:"name"`
}