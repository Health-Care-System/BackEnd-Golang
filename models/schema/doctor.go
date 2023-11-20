package schema

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID                uint                `gorm:"primarykey"`
	ProfilePicture    string              `gorm:"not null"`
	Fullname          string              `gorm:"not null"`
	Gender            string              `gorm:"type:enum('male', 'female')"`
	Email             string              `gorm:"not null;unique"`
	Password          string              `gorm:"not null"`
	Status            bool                `gorm:"not null;default:false"`
	Price             int                 `gorm:"not null"`
	Specialist        string              `gorm:"not null"`
	Experience        string              `gorm:"not null"`
	NoSTR             int                 `gorm:"not null"`
	Role              string              `gorm:"type:enum('doctor');default:'doctor'"`
	Alumnus           string              `gorm:"not null"`
	// DoctorTransaction []DoctorTransaction `gorm:"ForeignKey:DoctorID;references:ID"`
	UpdatedAt         time.Time
	CreatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
