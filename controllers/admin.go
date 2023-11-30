package controllers

import (
	"healthcare/configs"
	"healthcare/middlewares"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/response"
	"log"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Admin Login
func LoginAdminController(c echo.Context) error {
	var loginRequest web.AdminLoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input login data"))
	}

	if err := helper.ValidateStruct(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var admin schema.Admin
	if err := configs.DB.Where("email = ?", loginRequest.Email).First(&admin).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("email not registered"))
	}

	if err := configs.DB.Where("email = ? AND password = ?", loginRequest.Email, loginRequest.Password).First(&admin).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("incorrect email or password"))
	}

	token, err := middlewares.GenerateToken(admin.ID, admin.Email, admin.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to generate jwt"))
	}
	adminLoginResponse := response.ConvertToAdminLoginResponse(admin)
	adminLoginResponse.Token = token

	return c.JSON(http.StatusOK, helper.SuccessResponse("login successful", adminLoginResponse))
}

// Update Admin
func UpdateAdminController(c echo.Context) error {
	userID := c.Get("userID")

	var updatedAdmin web.AdminUpdateRequest

	if err := c.Bind(&updatedAdmin); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input update data"))
	}

	var existingAdmin schema.Admin
	result := configs.DB.First(&existingAdmin, userID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve admin"))
	}

	configs.DB.Model(&existingAdmin).Updates(updatedAdmin)

	response := response.ConvertToAdminUpdateResponse(&existingAdmin)

	return c.JSON(http.StatusOK, helper.SuccessResponse("admin updated data successful", response))
}

func AdminUpdatePaymentStatusController(c echo.Context) error {
	// Getting the admin ID from the context
	adminID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid admin id"))
	}

	// Log the admin ID for debugging purposes
	log.Printf("Admin ID: %d", adminID)

	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid transaction ID"))
	}

	var requestBody web.UpdatePaymentsRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid request body"))
	}

	// Checking if the required fields have been provided
	if requestBody.PaymentStatus == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("payment status is required"))
	}

	var doctorTransaction schema.DoctorTransaction
	err = configs.DB.First(&doctorTransaction, "id = ?", transactionID).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse("doctor transaction not found"))
	}

	// Updating payment status if provided
	doctorTransaction.PaymentStatus = requestBody.PaymentStatus

	// Saving the updated doctor transaction to the database
	if err := configs.DB.Save(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to update payment status"))
	}

	// Getting user data
	var user schema.User
	err = configs.DB.First(&user, "id=?", doctorTransaction.UserID).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user data"))
	}

	response := response.ConvertToPaymentsResponse(&doctorTransaction)
	return c.JSON(http.StatusOK, helper.SuccessResponse("payment status successfully updated", response))
}

// UpdatePaymentStatusByAdminController updates payment status by admin
// func UpdatePaymentStatusByAdminController(c echo.Context) error {
// 	// Parse transaction ID from the request parameters
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid transaction id"))
// 	}

// 	// Retrieve the existing transaction from the database
// 	var existingTransaction schema.DoctorTransaction
// 	result := configs.DB.First(&existingTransaction, id)
// 	if result.Error != nil {
// 		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve transaction"))
// 	}

// 	// Bind the updated payment status from the request body
// 	var updateRequest struct {
// 		PaymentStatus string `json:"payment_status" validate:"required,oneof=pending success cancelled"`
// 	}

// 	if err := c.Bind(&updateRequest); err != nil {
// 		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input update data"))
// 	}

// 	// Validate the updated payment status
// 	if err := helper.ValidateStruct(updateRequest); err != nil {
// 		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
// 	}

// 	return c.JSON(http.StatusOK, helper.SuccessResponse("payment status updated successfully", nil))
// }
