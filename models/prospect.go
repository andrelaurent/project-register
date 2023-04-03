package models

import "gorm.io/gorm"

type Prospect struct {
	gorm.Model
	ProspectID     string      `json:"ID" gorm:"primaryKey"`
	ProspectTypeID string      `json:"type_id" gorm:"index"`
	ProspectType   ProjectType `json:"prospect_type"`
	ProspectName   string      `json:"name"`
	UniqueNO       int         `json:"no"`
	Year           int         `json:"year"`
	ManagerID      string      `json:"manager_id"`
	Manager        Manager     `json:"manager"`
	ProspectStatus string      `json:"status"`
	ProspectTitle  string      `json:"title"`
	ProspectAmount int         `json:"amount"`
	CompanyID      string      `json:"company_id"`
	Company        Company     `json:"company"`
	ClientID       string      `json:"client_id"`
	Client         Client      `json:"client"`
	Jira           bool        `json:"jira"`
	Clockify       bool        `json:"clockify"`
	Pcs            bool        `json:"pcs"`
	Pms            bool        `json:"pms"`
}

type Prospects struct {
	Prospects []Prospect `json:"Prospects"`
}
