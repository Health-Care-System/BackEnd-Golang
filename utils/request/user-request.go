package request

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToUserRegisterRequest(user web.UserRegisterRequest) *schema.User {
	return &schema.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}
}

func ConvertToUserLoginRequest(user web.UserLoginRequest) *schema.User {
	return &schema.User{
		Email:    user.Email,
		Password: user.Password,
	}
}

func ConvertToUserUpdateRequest(user web.UserUpdateRequest) *schema.User {
	return &schema.User{
		Fullname:  user.Fullname,
		Email:     user.Email,
		Password:  user.Password,
		Image:     user.Image,
		Gender:    user.Gender,
		Birthdate: user.Birthdate,
		BloodType: user.BloodType,
		Height:    user.Height,
		Weight:    user.Weight,
	}
}