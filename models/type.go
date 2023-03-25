package models

import "gorm.io/gorm"

type ProjectType struct {
	gorm.Model
	ProjectTypeID   string `json:"code" gorm:"primaryKey;column:id"`
	ProjectTypeName string `json:"name"`
}
