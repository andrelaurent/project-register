package models

import "gorm.io/gorm"

type ProjectAudit struct {
	gorm.Model
	ID          uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProjectID   uint   `json:"project_id"`
	ProjectCode string `json:"code"`
	ProjectName string `json:"name"`
	Action      string `json:"action"`
	Date        string `json:"date"`
}

type ClientAudit struct {
	gorm.Model
	ID         uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientID   uint   `json:"client_id"`
	ClientCode string `json:"code"`
	ClientName string `json:"name"`
	Action     string `json:"action"`
	Date       string `json:"date"`
}

type CompanyAudit struct {
	gorm.Model
	ID          uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	CompanyID   uint   `json:"company_id"`
	CompanyCode string `json:"code"`
	CompanyName string `json:"name"`
	Action      string `json:"action"`
	Date        string `json:"date"`
	
}

type ProjectTypeAudit struct {
	gorm.Model
	ID              uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ProjectTypeID   uint   `json:"project_type_id"`
	ProjectTypeCode string `json:"code"`
	ProjectTypeName string `json:"name"`
	Action          string `json:"action"`
	Date            string `json:"date"`
}

type ContactAudit struct {
	gorm.Model
	ID          uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ContactID   uint   `json:"contact_id"`
	ContactName string `json:"name"`
	Action      string `json:"action"`
	Date        string `json:"date"`
	
}

type ClientContactAudit struct {
	gorm.Model
	ID              uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ClientContactID uint   `json:"client_contact_id"`
	ClientID        string `json:"client_id"`
	ContactID       string `json:"contact_id"`
	Action          string `json:"action"`
	Date            string `json:"date"`
}

type EmploymentAudit struct {
	gorm.Model
	ID              uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	EmploymentID    uint   `json:"employment_id"`
	ClientContactID uint   `json:"client_contact_id"`
	Action          string `json:"action"`
	Date            string `json:"date"`
}

type UserAudit struct {
	gorm.Model
	ID       uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Action   string `json:"action"`
	Date     string `json:"date"`
}

type LocationAudit struct {
	gorm.Model
	ID         uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	LocationID uint   `json:"location_id"`
	Address    string `json:"address"`
	Action     string `json:"action"`
	Date       string `json:"date"`
}
