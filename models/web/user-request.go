package web

type UserRegisterRequest struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required,max=30"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=10,max=15"`
}

type UserLoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=10,max=15"`
}

type UserUpdateRequest struct {
	Fullname  string `json:"fullname" form:"fullname" validate:"omitempty,max=30"`
	Email     string `json:"email" form:"email" validate:"omitempty,email"`
	Password  string `json:"password" form:"password" validate:"omitempty,min=10,max=15"`
	Image     string `json:"image" form:"image" validate:"omitempty"`
	Gender    string `json:"gender" form:"gender" validate:"omitempty"`
	Birthdate string `json:"birthdate" form:"birthdate" validate:"omitempty"`
	BloodType string `json:"blood_type" form:"blood_type" validate:"omitempty"`
	Height    int    `json:"height" form:"height" validate:"omitempty"`
	Weight    int    `json:"weight" form:"weight" validate:"omitempty"`
}
