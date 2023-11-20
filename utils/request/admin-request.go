package request

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToAdminLoginRequest(admin web.AdminLoginRequest) *schema.Admin {
	return &schema.Admin{
		Email:    admin.Email,
		Password: admin.Password,
	}
}

func ConvertToAdminUpdateRequest(admin web.AdminUpdateRequest) *schema.Admin {
	return &schema.Admin{
		Name:     admin.Name,
		Email:    admin.Email,
		Password: admin.Password,
	}
}
