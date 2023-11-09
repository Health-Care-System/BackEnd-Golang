package web

type UserRegisterResponse struct {
	Fullname  string `json:"fullname" form:"fullname"`
	Email     string `json:"email" form:"email"`
	ImageURL  string `json:"image_url" form:"image_url"`
	Gender    string `json:"gender" form:"gender"`
	Birthdate string `json:"birthdate" form:"birthdate"`
	BloodType string `json:"blood_type" form:"blood_type"`
	Height    int    `json:"height" form:"height"`
	Weight    int    `json:"weight" form:"weight"`
}

type UserLoginResponse struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Token    string `json:"token" form:"token"`
}

type UserUpdateResponse struct {
	Fullname  string `json:"fullname" form:"fullname"`
	Email     string `json:"email" form:"email"`
	Password  string `json:"password" form:"password"`
	ImageURL  string `json:"image_url" form:"image_url"`
	Gender    string `json:"gender" form:"gender"`
	Birthdate string `json:"birthdate" form:"birthdate"`
	BloodType string `json:"blood_type" form:"blood_type"`
	Height    int    `json:"height" form:"height"`
	Weight    int    `json:"weight" form:"weight"`
}
