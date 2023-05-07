package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID              uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
