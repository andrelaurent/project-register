package models

import "gorm.io/gorm"

type ProjectType struct {
	gorm.Model
	ProjectTypeCode string `json:"code" gorm:"primaryKey"`
	ProjectTypeName string `json:"name"`
}