package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID             uint           `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientCode     string         `json:"client_code"`
	ClientName     string         `json:"client_name"`
	Alias          string         `json:"alias"`
	Locations      []*Locations   `json:"locations" gorm:"foreignKey:ClientID"`
	SocialPresence SocialPresence `json:"client_social_presence" gorm:"embedded"`
	Subsidiary     Subsidiary     `json:"subsidiary" gorm:"embedded"`
	Date           time.Time      `json:"date"`
}

type SocialPresence struct {
	gorm.Model
	Linkedin string   `json:"linkedin"`
	Facebook string   `json:"facebook"`
	Twitter  string   `json:"twitter"`
	Github   string   `json:"github"`
	Other    []string `json:"other" gorm:"type:text[]"`
}

type Subsidiary struct {
	gorm.Model
	Subsidiaries     []string `json:"subsidiaries" gorm:"type:uuid[]"`
	ImmidiateParents []string `json:"immidiate_parents" gorm:"type:uuid[]"`
	UltimateParents  []string `json:"ultimate_parents" gorm:"type:uuid[]"`
}
