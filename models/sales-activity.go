package models

import "gorm.io/gorm"

type Activity struct {
	gorm.Model
	ID           uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ActivityName string `json:"activity_name"`
	ProjectId    uint   `json:"project_id"`
	Todos        []ToDo `json:"todos"`
}

type ToDo struct {
	gorm.Model
	ID         uint   `json:"ID" gorm:"primaryKey;autoIncrement"`
	ToDoName   string `json:"to_do_name"`
	ActivityId uint   `json:"activity_id"`
}
