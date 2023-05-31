package models

import (
	"gorm.io/gorm"
)

type ClientContact struct {
	gorm.Model
	ID          uint         `json:"ID" gorm:"primaryKey:autoIncrement"`
	ClientID    uint         `json:"client_id"`
	ContactID   uint         `json:"contact_id"`
	Employments []Employment `json:"employments"`
	ClientInfo  ClientInfo   `json:"client_info" gorm:"embedded"`
	ContactInfo ContactInfo  `json:"contact_info" gorm:"embedded"`
}

type Employment struct {
	gorm.Model
	JobTitle        string `json:"job_title"`
	JobStart        string `json:"job_start"`
	JobEnd          string `json:"job_end"`
	Status          string `json:"status"`
	ClientContactID uint   `json:"client_contact_id"`
}

type ClientInfo struct {
	ClientName string `json:"client_name"`
}

type ContactInfo struct {
	ContactName string `json:"contact_name"`
	BirthDate   string `json:"birth_date"`
}
