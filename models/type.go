package models

import "gorm.io/gorm"

type ProjectType struct {
	gorm.Model
	ID              uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProjectTypeCode string `json:"project_type_code"`
	ProjectTypeName string `json:"project_name"`
}
