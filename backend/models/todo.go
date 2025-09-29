package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id" gorm:"not null"`
	User        User   `json:"-" gorm:"foreignKey:UserID"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed" gorm:"default:false"`
}
