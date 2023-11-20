package web

type AdminLoginResponse struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}
type AdminUpdateResponse struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}
