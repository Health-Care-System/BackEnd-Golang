package controllers

import (
	"healthcare/configs"
	"healthcare/middlewares"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// User Register
func RegisterUserController(c echo.Context) error {
	var user web.UserRegisterRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Register Data"))
	}

	if err := helper.ValidateStruct(user); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	userRequest := request.ConvertToUserRegisterRequest(user)

	if existingUser := configs.DB.Where("email = ?", userRequest.Email).First(&userRequest).Error; existingUser == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("Email Already Exist"))
	}

	userRequest.Password = helper.HashPassword(userRequest.Password)

	if err := configs.DB.Create(&userRequest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Register"))
	}

	response := response.ConvertToUserRegisterResponse(userRequest)

	return c.JSON(http.StatusCreated, helper.SuccessResponse("Registered Successful", response))
}

// User Login
func LoginUserController(c echo.Context) error {
	var loginRequest web.UserLoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Login Data"))
	}

	if err := helper.ValidateStruct(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var user schema.User
	if err := configs.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Email Not Registered"))
	}

	if err := helper.ComparePassword(user.Password, loginRequest.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("Incorrect Email or Password"))
	}

	userLoginResponse := response.ConvertToUserLoginResponse(user)

	token, err := middlewares.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Generate JWT"))
	}

	userLoginResponse.Token = token

	return c.JSON(http.StatusOK, helper.SuccessResponse("Login Successful", userLoginResponse))
}

// Get User Profile
func GetUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Invalid User ID"))
	}

	var user schema.User

	if err := configs.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Retrieve User Data"))
	}

	response := response.ConvertToGetUserResponse(&user)

	return c.JSON(http.StatusOK, helper.SuccessResponse("Users Data Successfully Retrieved", response))
}

// Update User Profile
func UpdateUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Invalid User ID"))
	}

	var existingUser schema.User

	result := configs.DB.First(&existingUser, userID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Retrieve User"))
	}

	var userUpdated web.UserUpdateRequest

	if err := c.Bind(&userUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Update Data"))
	}

	if existingUser := configs.DB.Where("email = ?", userUpdated.Email).First(&userUpdated).Error; existingUser == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("Email Already Exist"))
	}

	if err := helper.ValidateStruct(userUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	userUpdated.Password = helper.HashPassword(userUpdated.Password)
	gender := strings.ToLower(userUpdated.Gender)
	bloodType := strings.ToUpper(userUpdated.BloodType)
	birthdate := userUpdated.Birthdate

	if !helper.GenderIsValid(gender) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Gender Data ('male', 'female')"))
	}

	if !helper.BloodTypeIsValid(bloodType) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Blood Type Data,('A', 'B', 'O', 'AB')"))
	}

	if !helper.BirthdateIsValid(birthdate) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid Input Birthdate Data (YYYY-MM-DD)"))
	}

	configs.DB.Model(&existingUser).Updates(userUpdated)

	userResponse := response.ConvertToUserUpdateResponse(&existingUser)

	return c.JSON(http.StatusOK, helper.SuccessResponse("User Updated Data Successful", userResponse))
}

// Delete User
func DeleteUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Invalid User ID"))
	}

	var existingUser schema.User
	result := configs.DB.First(&existingUser, userID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to Retrieve User"))
	}

	configs.DB.Delete(&existingUser)

	return c.JSON(http.StatusOK, helper.SuccessResponse("User Deleted Data Successful", nil))
}
