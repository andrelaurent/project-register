package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	CompanyID   string `json:"ID" gorm:"primaryKey;column:id;not null"`
	CompanyName string `json:"name" gorm:"not null"`
}
