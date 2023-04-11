package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	ID          uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	CompanyCode string `json:"company_code"`
	CompanyName string `json:"company_name"`
}
