package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/utils/helper"
	"healthcare/utils/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var db = configs.DB

func GetAvailableDoctor(c echo.Context) error {

	var doctors []schema.Doctor
	if err := db.Where("status = ?", true).Find(&doctors).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse("Data tidak ditemukan!"))
	}

	result := response.ConvertToAvailableDoctors(doctors)

	return c.JSON(http.StatusOK, helper.SuccessResponse("success", result))
}

func GetSpecializeDoctor(c echo.Context) error {
	specialist := c.QueryParam("specialist")

	if specialist == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Parameter specialist required!"))
	}

	var doctors []schema.Doctor
	if err := db.Where("specialist = ?", specialist).Find(&doctors).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("Data tidak ditemukan!"))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
	}

	result := response.ConvertToAvailableDoctors(doctors)

	return c.JSON(http.StatusOK, helper.SuccessResponse("success", result))
}
