package web

type UserRegisterResponse struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
}

type UserLoginResponse struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Token    string `json:"token" form:"token"`
}

type UserUpdateResponse struct {
	Fullname  string `json:"fullname" form:"fullname"`
	Email     string `json:"email" form:"email"`
	Image     string `json:"image" form:"image"`
	Gender    string `json:"gender" form:"gender"`
	Birthdate string `json:"birthdate" form:"birthdate"`
	BloodType string `json:"blood_type" form:"blood_type"`
	Height    int    `json:"height" form:"height"`
	Weight    int    `json:"weight" form:"weight"`
}
