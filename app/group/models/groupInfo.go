package models

import "gorm.io/gorm"

type GroupInfo struct {
	GroupID int `json:"group_id" gorm:"primary_key"`
	gorm.Model
}
