package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	Name    string `gorm:"size:255"`
	Content string `gorm:"size:text"`
	UserID  uint64 `gorm:"index"`
	User    User
}
