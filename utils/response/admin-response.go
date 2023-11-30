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

func ConvertToPaymentsResponse(admins *schema.DoctorTransaction) web.UpdatePaymentsResponse {
	return web.UpdatePaymentsResponse{
		TransactionID: admins.ID,
		PaymentStatus: admins.PaymentStatus,
	}
}
