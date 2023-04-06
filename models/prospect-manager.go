package models

import "gorm.io/gorm"

type ProspectManager struct {
	gorm.Model
	ManagerID   string `json:"ID" gorm:"primaryKey;column:id;not null"`
	ManagerName string `json:"name" gorm:"not null"`
}
