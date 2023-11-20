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
