package schema

import (
	"time"
)

type Admin struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	UpdatedAt time.Time
}
