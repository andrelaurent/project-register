package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	ID                    uint           `json:"ID" gorm:"primaryKey;autoIncrement"`
	ContactName           string         `json:"contact_name"`
	ContactAlias          string         `json:"contact_alias"`
	Gender                string         `json:"gender"`
	Emails                pq.StringArray `json:"contact_emails" gorm:"type:text[]"`
	Phones                pq.StringArray `json:"contact_phones" gorm:"type:text[]"`
	ContactSocialPresence SocialPresence `json:"contact_social_presence" gorm:"embedded"`
	BirthDate             string         `json:"birth_date"`
	Religion              string         `json:"religion"`
	Interests             pq.StringArray `json:"interests" gorm:"type:text[]"`
	Skills                pq.StringArray `json:"skills" gorm:"type:text[]"`
	Educations            pq.StringArray `json:"educations" gorm:"type:text[]"`
	Notes                 string         `json:"notes"`
	Locations             []Locations    `json:"locations"`
}
