package web

type MedicineRequest struct {
	Code     string `json:"code" form:"code" validate:"required"`
	Name     string `json:"name" form:"name" validate:"required"`
	Merk     string `json:"merk" form:"merk" validate:"required"`
	Category string `json:"category" form:"category" validate:"required"`
	Type     string `json:"type" form:"type" validate:"required"`
	Stock    int    `json:"stock" form:"stock" validate:"required,min=0"`
	Price    int    `json:"price" form:"price" validate:"required,min=0"`
	Details  string `json:"details" form:"details" validate:"required"`
	Image    string `json:"image" form:"image" validate:"required"`
}
