package schema

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID                 uint   `gorm:"primarykey"`
	Fullname           string `gorm:"not null"`
	Email              string `gorm:"not null;unique"`
	Password           string `gorm:"not null"`
	Status             string `gorm:"not null"`
	Price              int    `gorm:"not null"`
	Tag                string `gorm:"not null"`
	ProfilePicture     string `gorm:"not null"`
	RegistrationLetter string `gorm:"not null"`
	Role               string `gorm:"type:enum('doctor');default:'doctor'"`
	UpdatedAt          time.Time
	CreatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
