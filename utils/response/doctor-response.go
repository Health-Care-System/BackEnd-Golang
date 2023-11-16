package response

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToDoctorRegisterResponse(doctor *schema.Doctor) web.DoctorResgisterResponse {
	return web.DoctorResgisterResponse{
		Fullname:           doctor.Fullname,
		Email:              doctor.Email,
		Price:              doctor.Price,
		Tag:                doctor.Tag,
		ProfilePicture:     doctor.ProfilePicture,
		RegistrationLetter: doctor.RegistrationLetter,
	}
}

func ConvertToDoctorLoginResponse(doctor *schema.Doctor) web.DoctorLoginResponse {
	return web.DoctorLoginResponse{
		Fullname: doctor.Fullname,
		Email:    doctor.Email,
	}
}

func ConvertToDoctorUpdateResponse(doctor *schema.Doctor) web.DoctorUpdateResponse {
	return web.DoctorUpdateResponse{
		Fullname:           doctor.Fullname,
		Email:              doctor.Email,
		Price:              doctor.Price,
		Tag:                doctor.Tag,
		ProfilePicture:     doctor.ProfilePicture,
		RegistrationLetter: doctor.RegistrationLetter,
	}
}

func ConvertToGetDoctorResponse(doctor *schema.Doctor) web.DoctorResgisterResponse {
	return web.DoctorResgisterResponse{
		Fullname:           doctor.Fullname,
		Email:              doctor.Email,
		Price:              doctor.Price,
		Tag:                doctor.Tag,
		ProfilePicture:     doctor.ProfilePicture,
		RegistrationLetter: doctor.RegistrationLetter,
	}
}

func ConvertToGetAllDoctorResponse(doctors []schema.Doctor) []web.DoctorAllResponse {
	var results []web.DoctorAllResponse

	// Iterasi melalui setiap dokter dan konversi ke format respons
	for _, doctor := range doctors {
		doctorResponse := web.DoctorAllResponse{
			Fullname:           doctor.Fullname,
			Email:              doctor.Email,
			Price:              doctor.Price,
			Tag:                doctor.Tag,
			Status:             doctor.Status,
			ProfilePicture:     doctor.ProfilePicture,
			RegistrationLetter: doctor.RegistrationLetter,
		}

		results = append(results, doctorResponse)
	}

	return results
}

// func ConvertToSpecializeDoctorResponse(doctor schema.Doctor) web.SpecializeDoctorResponse {
// 	return web.SpecializeDoctorResponse{
// 		ID:       doctor.ID,
// 		Fullname: doctor.Fullname,
// 		Email:    doctor.Email,
// 		Status:   doctor.Status,
// 		Price:    doctor.Price,
// 		Tag:      doctor.Tag,
// 		ImageURL: doctor.ProfilePicture,
// 	}
// }
