package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientID   string `json:"ID" gorm:"primaryKey;column:id;not null"`
	ClientName string `json:"name" gorm:"not null"`
}
