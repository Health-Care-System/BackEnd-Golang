package response

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToAvailableDoctors(doctors []schema.Doctor) []web.AvailableDoctorResponse {
	var results []web.AvailableDoctorResponse

	for _, doctor := range doctors {
		doctorResponses := web.AvailableDoctorResponse{
			ID:             doctor.ID,
			ProfilePicture: doctor.ProfilePicture,
			Fullname:       doctor.Fullname,
			Gender:         doctor.Gender,
			Status:         doctor.Status,
			Price:          doctor.Price,
			Specialist:     doctor.Specialist,
			Experience:     doctor.Experience,
			NoSTR:          doctor.NoSTR,
			Alumnus:        doctor.Alumnus,
		}
		results = append(results, doctorResponses)
	}
	return results
}

func ConvertToDoctorLoginResponse(doctor *schema.Doctor) web.DoctorLoginResponse {
	return web.DoctorLoginResponse{
		Fullname: doctor.Fullname,
		Email:    doctor.Email,
	}
}

func ConvertToDoctorUpdateResponse(doctor *schema.Doctor) web.DoctorUpdateResponse {
	return web.DoctorUpdateResponse{
		Fullname:       doctor.Fullname,
		Email:          doctor.Email,
		Gender:         doctor.Gender,
		Specialist:     doctor.Specialist,
		ProfilePicture: doctor.ProfilePicture,
		NoSTR:          doctor.NoSTR,
		Status:         doctor.Status,
		Experience:     doctor.Experience,
		Alumnus:        doctor.Alumnus,
	}
}

func ConvertToGetDoctorResponse(doctor *schema.Doctor) web.DoctorUpdateResponse {
	return web.DoctorUpdateResponse{
		Fullname:       doctor.Fullname,
		Email:          doctor.Email,
		Gender:         doctor.Gender,
		Specialist:     doctor.Specialist,
		ProfilePicture: doctor.ProfilePicture,
		NoSTR:          doctor.NoSTR,
		Experience:     doctor.Experience,
		Alumnus:        doctor.Alumnus,
	}
}

func ConvertToGetAllDoctorResponse(doctors []schema.Doctor) []web.DoctorAllResponse {
	var results []web.DoctorAllResponse

	for _, doctor := range doctors {
		doctorResponse := web.DoctorAllResponse{
			ProfilePicture: doctor.ProfilePicture,
			Fullname:       doctor.Fullname,
			NoSTR:          doctor.NoSTR,
			Gender:         doctor.Gender,
			Status:         doctor.Status,
			Price:          doctor.Price,
			Specialist:     doctor.Specialist,
			Experience:     doctor.Experience,
			Alumnus:        doctor.Alumnus,
		}

		results = append(results, doctorResponse)
	}

	return results
}