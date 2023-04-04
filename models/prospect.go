package models

import "gorm.io/gorm"

type Prospect struct {
	gorm.Model
	ProspectID     string      `json:"ID" gorm:"primaryKey;column:id;not null"`
	ProjectTypeID  string      `json:"type_id" gorm:"index;not null"`
	ProjectType    ProjectType `json:"project_type" gorm:"not null"`
	ProspectName   string      `json:"name" gorm:"not null"`
	UniqueNO       int         `json:"no" gorm:"not null"`
	Year           int         `json:"year" gorm:"not null"`
	ManagerID      string      `json:"manager_id" gorm:"index;not null"`
	Manager        Manager     `json:"manager" gorm:"not null"`
	ProspectStatus string      `json:"status" gorm:"not null"`
	ProspectTitle  string      `json:"title" gorm:"not null"`
	ProspectAmount int         `json:"amount" gorm:"not null"`
	CompanyID      string      `json:"company_id" gorm:"index;not null"`
	Company        Company     `json:"company" gorm:"not null"`
	ClientID       string      `json:"client_id" gorm:"index;not null"`
	Client         Client      `json:"client" gorm:"not null"`
	Jira           bool        `json:"jira" gorm:"not null"`
	Clockify       bool        `json:"clockify" gorm:"not null"`
	Pcs            bool        `json:"pcs" gorm:"not null"`
	Pms            bool        `json:"pms" gorm:"not null"`
}

type Prospects struct {
	Prospects []Prospect `json:"Prospects"`
}
