package schema

import "time"

type Message struct {
	ID         uint `gorm:"primaryKey"`
	RoomchatID uint 
	UserID     uint
	DoctorID   uint
	Message    string
	Image      string
	Audio      string
	CreatedAt  time.Time
}



