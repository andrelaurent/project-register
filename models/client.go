package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID             uuid.UUID      `json:"ID" gorm:"type:uuid;primaryKey"`
	ClientCode     string         `json:"client_code"`
	ClientName     string         `json:"client_name"`
	Alias          string         `json:"alias"`
	Locations      []Locations    `json:"locations"`
	SocialPresence SocialPresence `json:"client_social_presence"`
	Subsidiary     Subsidiary     `json:"subsidiary"`
	Date           time.Time      `json:"date"`
}

type SocialPresence struct {
	gorm.Model
	Linkedin string   `json:"linkedin"`
	Facebook string   `json:"facebook"`
	Twitter  string   `json:"twitter"`
	Github   string   `json:"github"`
	Other    []string `json:"other"`
}

type Subsidiary struct {
	gorm.Model
	Subsidiaries     []uuid.UUID `gorm:"type:uuid[]"`
	ImmidiateParents []uuid.UUID `gorm:"type:uuid[]"`
	UltimateParents  []uuid.UUID `gorm:"type:uuid[]"`
}
