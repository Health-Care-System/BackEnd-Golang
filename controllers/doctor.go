package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAvailableDoctor(c echo.Context) error {
	db := configs.DB

	var doctors []schema.Doctor
	if err := db.Where("status = ?", "true").Find(&doctors).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse("Data tidak ditemukan!"))
	}

	var doctorResponses []web.AvailableDoctorResponse

	for _, doctor := range doctors {
		doctorResponses = append(doctorResponses, web.AvailableDoctorResponse{
			ID:       int(doctor.ID),
			Fullname: doctor.Fullname,
			Email:    doctor.Email,
			Price:    doctor.Price,
			Tag:      doctor.Tag,
			ImageURL: doctor.ProfilePicture,
		})
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("success", doctorResponses))
}

func GetSpecializeDoctor(c echo.Context) error {
	tag := c.QueryParam("tag")

	if tag == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Parameter tag required"))
	}

	var doctors []schema.Doctor
	if err := configs.DB.Where("tag = ?", tag).Find(&doctors).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("Data tidak ditemukan"))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
	}

	var specializeDoctors []web.SpecializeDoctorResponse
	for _, doctor := range doctors {
		specializeDoctor := web.SpecializeDoctorResponse{
			ID:       doctor.ID,
			Fullname: doctor.Fullname,
			Email:    doctor.Email,
			Status:   doctor.Status,
			Price:    doctor.Price,
			Tag:      doctor.Tag,
			ImageURL: doctor.ProfilePicture,
		}
		specializeDoctors = append(specializeDoctors, specializeDoctor)
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("success", specializeDoctors))
}
