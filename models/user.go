package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}
