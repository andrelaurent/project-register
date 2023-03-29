package models

import "gorm.io/gorm"

type Manager struct {
	gorm.Model
	ManagerID   string `json:"ID" gorm:"primaryKey;column:id"`
	ManagerName string `json:"name"`
}
