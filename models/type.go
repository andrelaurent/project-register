package models

import "gorm.io/gorm"

type ProjectType struct {
	gorm.Model
	ProjectTypeID   string `json:"ID" gorm:"primaryKey;column:id;not null;default:null"`
	ProjectTypeName string `json:"name" gorm:"not null;default:null"`
}
