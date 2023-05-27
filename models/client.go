package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID             uint           `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientCode     string         `json:"client_code"`
	ClientName     string         `json:"client_name"`
	Alias          string         `json:"alias"`
	LocationIDs    []uint         `json:"location_ids" gorm:"type:integer[]"`
	Website        string         `json:"website"`
	SocialPresence SocialPresence `json:"client_social_presence" gorm:"embedded"`
	Subsidiary     Subsidiary     `json:"subsidiary" gorm:"embedded type:jsonb"`
	Date           time.Time      `json:"date"`
}

type SocialPresence struct {
	Linkedin string   `json:"linkedin"`
	Facebook string   `json:"facebook"`
	Twitter  string   `json:"twitter"`
	Github   string   `json:"github"`
	Other    []string `json:"other" gorm:"type:text[]"`
}

type Subsidiary struct {
	gorm.Model
	Subsidiaries     pq.StringArray `json:"subsidiaries" gorm:"type:text[]"`
	ImmidiateParents pq.StringArray `json:"immidiate_parents" gorm:"type:text[]"`
	UltimateParents  pq.StringArray `json:"ultimate_parents" gorm:"type:text[]"`
}
