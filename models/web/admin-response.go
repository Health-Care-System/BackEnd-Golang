package web

type AdminLoginResponse struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Token string `json:"token" form:"token"`
}
type AdminUpdateResponse struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

type AdminProfileResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
