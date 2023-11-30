package response

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToAdminLoginResponse(admin schema.Admin) web.AdminLoginResponse {
	return web.AdminLoginResponse{
		Name:  admin.Name,
		Email: admin.Email,
	}
}
func ConvertToAdminUpdateResponse(admin *schema.Admin) web.AdminUpdateResponse {
	return web.AdminUpdateResponse{
		Name:  admin.Name,
		Email: admin.Email,
	}
}

func ConvertToGetProfileAdminResponse(admin *schema.Admin) web.AdminProfileResponse {
	return web.AdminProfileResponse{
		Name:  admin.Name,
		Email: admin.Email,
	}
}
