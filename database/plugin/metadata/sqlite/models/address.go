package models

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	Address   string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Address) TableName() string {
	return "Address"
}
