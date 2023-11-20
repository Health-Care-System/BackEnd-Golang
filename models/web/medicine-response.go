package web

import "time"

type MedicineResponse struct {
	ID        uint      `json:"id" form:"id"`
	Code      string    `json:"code" form:"code"`
	Name      string    `json:"name" form:"name"`
	Merk      string    `json:"merk" form:"merk"`
	Category  string    `json:"category" form:"category"`
	Type      string    `json:"type" form:"type"`
	Price     int       `json:"price" form:"price"`
	Stock     int       `json:"stock" form:"stock"`
	Details   string    `json:"details" form:"details"`
	Image     string    `json:"image" form:"image"`
	CreatedAt time.Time `json:"created_at" form:"created_at"`
}
