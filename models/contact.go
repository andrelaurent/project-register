package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID                    uint           `json:"ID" gorm:"primaryKey;autoIncrement"`
	ContactName           string         `json:"contact_name"`
	ContactAlias          string         `json:"contact_alias"`
	Gender                byte           `json:"gender"`
	Emails                []string       `json:"contact_emails" gorm:"type:text[]"`
	Phones                []string       `json:"contact_phones" gorm:"type:text[]"`
	Locations             []Locations    `json:"contact_locations" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ContactSocialPresence SocialPresence `json:"contact_social_presence" gorm:"embedded"`
	BirthDate             time.Time      `json:"birth_date"`
	Religion              string         `json:"religion"`
	Interests             []string       `json:"interests" gorm:"type:text[]"`
	Skills                []string       `json:"skills" gorm:"type:text[]"`
	Educations            []string       `json:"educations" gorm:"type:text[]"`
	Notes                 string         `json:"notes"`
}
