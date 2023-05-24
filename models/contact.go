package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID uint `json:"ID" gorm:"primaryKey;autoIncrement"`
	ContactName string `json:"contact_name"`
	ContactAlias string `json:"contact_alias"`
	Gender byte `json:"gender"`
	Emails []string `json:"contact_emails"`
	Phones []string `json:"contact_phones"`
	Locations []Locations `json:"contact_locations"`
	SocialPresence SocialPresence
}