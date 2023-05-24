package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID uuid.UUID `json:"ID" gorm:"type:uuid;primaryKey"`
	ContactName string `json:"contact_name"`
	ContactAlias string `json:"contact_alias"`
	Gender byte `json:"gender"`
	Emails []string `json:"contact_emails"`
	Phones []string `json:"contact_phones"`
	Locations []Locations `json:"contact_locations"`
	ContactSocialPresence SocialPresence `json:"contact_social_presence"`
	BirthDate time.Time `json:"birth_date"`
	Religion string `json:"religion"`
	Interests []string `json:"interests"`
	Skills []string `json:"skills"`
	Educations []string `json:"educations"`
	Notes string `json:"notes"`
}