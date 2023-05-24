package models

import (
	"github.com/andrelaurent/project-register/vendor/github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Client struct {
	gorm.Model
	ID         uint        `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientCode string      `json:"client_code"`
	ClientName string      `json:"client_name"`
	Alias      string      `json:"alias"`
	Locations  []Locations `json:"locations"`
	Subsidiary Subsidiary  `json:"subsidiary"`
	Date       time.Time   `json:"date"`
}

type SocialPresence struct {
	gorm.Model
	Linkedin string `json:"linkedin"`
	Facebook string `json:"facebook"`
	Twitter  string `json:"Twitter"`
	Other    string `json:"other"`
}

type Subsidiary struct {
	gorm.Model
	Subsidiaries     []uuid.UUID `gorm:"type:uuid[]"`
	ImmidiateParents []uuid.UUID `gorm:"type:uuid[]"`
	UltimateParents  []uuid.UUID `gorm:"type:uuid[]"`
}