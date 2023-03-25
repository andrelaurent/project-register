package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyID string `json:"code" gorm:"primaryKey;column:id"`
	CompanyName string `json:"name"`
}