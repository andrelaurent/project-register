package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID         uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientCode string `json:"client_code"`
	ClientName string `json:"client_name"`
}
