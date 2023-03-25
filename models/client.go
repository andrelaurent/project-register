package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ClientID string `json:"code" gorm:"primaryKey"`
	ClientName string `json:"name"`
}
