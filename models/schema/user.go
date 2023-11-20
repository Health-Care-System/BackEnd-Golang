package schema

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	Fullname  string
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Image     string
	Gender    string `gorm:"type:enum('male', 'female', '')"`
	Birthdate string 
	BloodType string `gorm:"type:enum('A', 'B', 'O', 'AB', '')"`
	Height    int
	Weight    int
	Role      string `gorm:"type:enum('user');default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
