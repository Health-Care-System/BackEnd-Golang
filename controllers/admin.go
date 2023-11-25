package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/response"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Admin Login
func LoginAdminController(c echo.Context) error {
	var loginRequest web.AdminLoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Login Data"))
	}

	if err := helper.ValidateStruct(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var admin schema.Admin
	if err := configs.DB.Where("email = ?", loginRequest.Email).First(&admin).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Email Not Registered"))
	}

	if err := configs.DB.Where("email = ? AND password = ?", loginRequest.Email, loginRequest.Password).First(&admin).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Incorrect Email or Password"))
	}

	adminLoginResponse := response.ConvertToAdminLoginResponse(admin)

	return c.JSON(http.StatusOK, helper.SuccessResponse("Login Successful", adminLoginResponse))
}

// Update Admin
func UpdateAdminController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Admin ID"))
	}

	var updatedAdmin web.AdminUpdateRequest

	if err := c.Bind(&updatedAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Update Data"))
	}

	var existingAdmin schema.Admin
	result := configs.DB.First(&existingAdmin, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Retrieve Admin"))
	}

	configs.DB.Model(&existingAdmin).Updates(updatedAdmin)

	response := response.ConvertToAdminUpdateResponse(&existingAdmin)

	return c.JSON(http.StatusOK, helper.SuccessResponse("Admin Updated Data Successful", response))
}

// UpdatePaymentStatusByAdminController updates payment status by admin
func UpdatePaymentStatusByAdminController(c echo.Context) error {
	// Parse transaction ID from the request parameters
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Transaction ID"))
	}

	// Retrieve the existing transaction from the database
	var existingTransaction schema.DoctorTransaction
	result := configs.DB.First(&existingTransaction, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Retrieve Transaction"))
	}

	// Bind the updated payment status from the request body
	var updateRequest struct {
		PaymentStatus string `json:"payment_status" validate:"required,oneof=pending success cancelled"`
	}

	if err := c.Bind(&updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Update Data"))
	}

	// Validate the updated payment status
	if err := helper.ValidateStruct(updateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}
	// Update the payment status in the existing transaction
	// existingTransaction.PaymentStatus = updateRequest.PaymentStatus
	// result = configs.DB.Save(&existingTransaction)
	// if result.Error != nil {
	// 	return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Update Payment Status"))
	// }

	return c.JSON(http.StatusOK, helper.SuccessResponse("Admin Updated Payment Status Successfully", nil))
}
