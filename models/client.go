package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientID string `json:"ID" gorm:"primaryKey;column:id"`
	ClientName string `json:"name"`
}
