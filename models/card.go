package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	Number string
	UserID uint64
}
