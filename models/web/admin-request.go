package web

type AdminLoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=10,max=15"`
}

type AdminUpdateRequest struct {
	Name     string `json:"name" form:"name" validate:"omitempty"`
	Email    string `json:"email" form:"email" validate:"omitempty,email"`
	Password string `json:"password" form:"password" validate:"omitempty,min=10,max=15"`
}

type UpdatePaymentsRequest struct {
	PaymentStatus string `json:"payment_status" form:"payment_status" validate:"required"`
}
