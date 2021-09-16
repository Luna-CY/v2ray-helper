package model

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	Type  string `gorm:"not null"`
	Name  string `gorm:"not null"`
	Value string `gorm:"not null"`
}

func (s *Setting) TableName() string {
	return "settings"
}
