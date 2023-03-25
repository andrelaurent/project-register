package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientID string `json:"code" gorm:"primaryKey;column:id"`
	ClientName string `json:"name"`
}
