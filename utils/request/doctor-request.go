package request

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToDoctorLoginRequest(doctor web.DoctorLoginRequest) *schema.Doctor {
	return &schema.Doctor{
		Email:    doctor.Email,
		Password: doctor.Password,
	}
}

func ConvertToDoctorUpdateRequest(doctor web.DoctorUpdateRequest) *schema.Doctor {
	return &schema.Doctor{
		Fullname:       doctor.Fullname,
		Email:          doctor.Email,
		Specialist:     doctor.Specialist,
		ProfilePicture: doctor.ProfilePicture,
		Gender:         doctor.Gender,
		Status:         doctor.Status,
		Experience:     doctor.Experience,
		Alumnus:        doctor.Alumnus,
	}
}
