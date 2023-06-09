package models

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID              uint        `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProjectID       string      `json:"project_id" gorm:"not null"`
	ProjectTypeID   uint        `json:"type_id" gorm:"index;not null"`
	ProjectType     ProjectType `json:"project_type" gorm:"not null;"`
	ProjectName     string      `json:"project_name" gorm:"not null"`
	UniqueNO        int         `json:"no" gorm:"not null"`
	Year            int         `json:"year" gorm:"not null"`
	Pic             string      `json:"manager" gorm:"not null"`
	ProjectStatus   string      `json:"status" gorm:"not null"`
	ProjectTitle    string      `json:"title" gorm:"not null"`
	ProjectAmount   float64     `json:"amount" gorm:"not null"`
	CompanyID       uint        `json:"company_id" gorm:"index;not null"`
	Company         Company     `json:"company" gorm:"not null;"`
	ClientID        uint        `json:"client_id" gorm:"index;not null"`
	Client          Client      `json:"client" gorm:"not null;"`
	IsDeleted       bool        `json:"is_deleted" gorm:"not null"`
	Jira            bool        `json:"jira" gorm:"not null"`
	Clockify        bool        `json:"clockify" gorm:"not null"`
	Pcs             bool        `json:"pcs" gorm:"not null"`
	Pms             bool        `json:"pms" gorm:"not null"`
	SalesActivities []Activity  `json:"sales_activities"`
}