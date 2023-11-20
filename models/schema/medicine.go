package schema

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	ID        uint   `gorm:"primarykey"`
	Code      string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Merk      string `gorm:"not null"`
	Category  string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Stock     int    `gorm:"not null"`
	Price     int    `gorm:"not null"`
	Details   string `gorm:"not null"`
	Image     string `gorm:"not null"`
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
