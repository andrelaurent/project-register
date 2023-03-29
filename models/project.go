package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ProjectID     string      `json:"ID" gorm:"primaryKey"`
	ProjectTypeID string      `json:"type_id" gorm:"index"`
	ProjectType   ProjectType `json:"project_type"`
	ProjectName   string      `json:"name"`
	UniqueNO      int         `json:"no"`
	Year          int         `json:"year"`
	ManagerID     string      `json:"manager_id"`
	Manager       Manager     `json:"manager"`
	ProjectStatus string      `json:"status"`
	ProjectTitle  string      `json:"title"`
	ProjectAmount int         `json:"amount"`
	CompanyID     string      `json:"company_id"`
	Company       Company     `json:"company"`
	ClientID      string      `json:"client_id"`
	Client        Client      `json:"client"`
	Jira          bool        `json:"jira"`
	Clockify      bool        `json:"clockify"`
	Pcs           bool        `json:"pcs"`
	Pms           bool        `json:"pms"`
}

type Projects struct {
	Projects []Project `json:"projects"`
}
