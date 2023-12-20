package controllers

import (
	"fmt"
	"healthcare/configs"
	"healthcare/middlewares"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/helper/constanta"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

// RegisterUserController
func RegisterUserController(c echo.Context) error {
	var user web.UserRegisterRequest

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(user); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	userRequest := request.ConvertToUserRegisterRequest(user)

	var existingUserEmail schema.User
	if existingEmail := configs.DB.Where("email = ? AND deleted_at IS NULL", userRequest.Email).First(&existingUserEmail).Error; existingEmail == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("email already exists"))
	}

	userRequest.Password = helper.HashPassword(userRequest.Password)
	if err := configs.DB.Create(&userRequest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"register"))
	}

	// Send OTP via email
	err := helper.SendOTPViaEmail(userRequest.Email, "user", "register")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"OTP via email"))
	}

	// Use the same success response for both cases
	response := response.ConvertToUserRegisterResponse(userRequest)
	return c.JSON(http.StatusCreated, helper.SuccessResponse(constanta.SuccessActionCreated+"registeredCheck your email for OTP verification", response))
}

// LoginUserController
func LoginUserController(c echo.Context) error {
	var loginRequest web.UserLoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input login data"))
	}

	if err := helper.ValidateStruct(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var user schema.User
	if err := configs.DB.Where("email = ? AND deleted_at IS NULL", loginRequest.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("email not registered"))
	}

	// Periksa apakah pengguna diverifikasi dengan OTP
	if !user.IsVerified {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("user is not verified"))
	}
	if err := helper.ComparePassword(user.Password, loginRequest.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("incorrect email or password"))
	}

	userLoginResponse := response.ConvertToUserLoginResponse(user)

	token, err := middlewares.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to generate jwt"))
	}

	userLoginResponse.Token = token

	err = helper.SendNotificationEmail(user.Email, user.Fullname, "login", "user", "", "", false, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to send verification email"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("login successful", userLoginResponse))
}

func GetAllUserPagination(offset int, limit int, queryInput []schema.User) ([]schema.User, int64, error) {

	if offset < 0 || limit < 0 {
		return nil, 0, nil
	}

	queryAll := queryInput
	var total int64

	query := configs.DB.Model(&queryAll)

	query.Find(&queryAll).Count(&total)

	query = query.Limit(limit).Offset(offset)

	result := query.Find(&queryAll)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	if offset >= int(total) {
		return nil, 0, fmt.Errorf("not found")
	}

	return queryAll, total, nil
}

// Get All Doctors by Admin
func GetAllUserByAdminController(c echo.Context) error {

	params := c.QueryParams()
	limit, err := strconv.Atoi(params.Get("limit"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(params.Get("offset"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	var users []schema.User
	user, total, err := GetAllUserPagination(offset, limit, users)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	pagination := helper.Pagination(offset, limit, total)

	response := response.ConvertToGetAllUserByAdminResponse(user)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"user", response, pagination))
}

// Get User Profile
func GetUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	var user schema.User

	if err := configs.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user data"))
	}

	response := response.ConvertToGetUserResponse(&user)

	return c.JSON(http.StatusOK, helper.SuccessResponse("users data successfully retrieved", response))
}

// Get User by ID
func GetUserIDbyAdminController(c echo.Context) error {
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid user id"))
	}
	var user schema.User
	result := configs.DB.First(&user, user_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user data"))
	}
	response := response.ConvertToGetUserIDbyAdminResponse(&user)

	return c.JSON(http.StatusOK, helper.SuccessResponse("users data successfully retrieved", response))
}

// Get User Transaction by Admin
func GetUserPaymentsByAdminsController(c echo.Context) error {
	// Mendapatkan ID pengguna dari URL
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid user_id"))
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	var doctorTransactions []schema.DoctorTransaction
	var total int64

	// Fetch transactions for the specified user with payment_status IN ('success', 'pending', 'cancelled')
	query := configs.DB.Where("user_id = ? AND payment_status IN (?)", userID, []string{"success", "pending", "cancelled"}).Find(&doctorTransactions)
	if query.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor transactions"))
	}

	sort.Slice(doctorTransactions, func(i, j int) bool {
		order := map[string]int{"pending": 0, "success": 1, "cancelled": 2}
		return order[doctorTransactions[i].PaymentStatus] < order[doctorTransactions[j].PaymentStatus]
	})

	// Count total number of records
	total = int64(len(doctorTransactions))

	// Apply limit and offset to the result
	start := offset
	end := offset + limit
	if start > len(doctorTransactions) {
		start = len(doctorTransactions)
	}
	if end > len(doctorTransactions) {
		end = len(doctorTransactions)
	}
	doctorTransactions = doctorTransactions[start:end]

	if len(doctorTransactions) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound+"doctor transactions"))
	}

	pagination := helper.Pagination(offset, limit, total)
	responses := response.ConvertToAdminDoctorPaymentsResponse(doctorTransactions)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionCreated+"doctor transactions", responses, pagination))
}

// Update User Profile
func UpdateUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		log.Println(ok)
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	var existingUser schema.User

	result := configs.DB.Where("id = ?", userID).First(&existingUser)
	if result.Error != nil {
		log.Println(result)
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user"))
	}

	var userUpdated web.UserUpdateRequest

	if err := c.Bind(&userUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input update data"))
	}

	var existingUserEmail schema.User
	if existingEmail := configs.DB.Where("email = ? AND deleted_at IS NULL", userUpdated.Email).First(&existingUserEmail).Error; existingEmail == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("email already exist"))
	}

	if err := helper.ValidateStruct(userUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("profile_picture")

	if err == nil {
		defer file.Close()

		if fileHeader.Size > 10*1024*1024 { // 10 MB limit
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("image file size exceeds the limit (10 MB)"))
		}

		allowedExtensions := []string{".jpg", ".jpeg", ".png"}
		ext := filepath.Ext(fileHeader.Filename)
		allowed := false
		for _, validExt := range allowedExtensions {
			if ext == validExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid image file format. supported formats: .jpg, .jpeg, .png"))
		}

		profilePicture, err := helper.UploadFilesToGCS(c, fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error uploading image to cloud storage"))
		}

		userUpdated.ProfilePicture = profilePicture
	}

	userUpdated.Password = helper.HashPassword(userUpdated.Password)
	gender := strings.ToLower(userUpdated.Gender)
	bloodType := strings.ToUpper(userUpdated.BloodType)
	birthdate := userUpdated.Birthdate

	if !helper.GenderIsValid(gender) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input gender data ('male', 'female')"))
	}

	if !helper.BloodTypeIsValid(bloodType) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input blood type data,('a', 'b', 'o', 'ab')"))
	}

	if !helper.BirthdateIsValid(birthdate) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input birthdate data (yyyy-mm-dd)"))
	}

	configs.DB.Model(&existingUser).Updates(userUpdated)

	userResponse := response.ConvertToUserUpdateResponse(&existingUser)

	return c.JSON(http.StatusOK, helper.SuccessResponse("user updated data successful", userResponse))
}

// Delete User
func DeleteUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	var existingUser schema.User
	result := configs.DB.First(&existingUser, userID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user"))
	}

	configs.DB.Delete(&existingUser)

	return c.JSON(http.StatusOK, helper.SuccessResponse("user deleted data successful", nil))
}

func DeleteUserByAdminController(c echo.Context) error {
	// Parse doctor ID from the request parameters
	user_id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid user id"))
	}

	// Retrieve the existing doctor from the database
	var existingUser schema.User
	result := configs.DB.First(&existingUser, user_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve user"))
	}

	// Delete the doctor from the database
	result = configs.DB.Delete(&existingUser, user_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to delete user"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("user deleted data successful  ", nil))
}

// VerifyOTP
func VerifyOTPRegister(c echo.Context) error {
	var verificationRequest web.OTPVerificationRequest
	if err := c.Bind(&verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid request"))
	}

	if err := helper.ValidateStruct(verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var wg sync.WaitGroup
	defer wg.Wait() // Pastikan menunggu goroutine selesai sebelum fungsi selesai

	if err := helper.VerifyOTPByEmail(verificationRequest.Email, verificationRequest.OTP, "user"); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionGet+"OTP not found"))
	}

	// Update user verification status
	if err := helper.UpdateUserVerificationStatus(verificationRequest.Email, true); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+" update user verification status"))
	}

	// Send registration notification email
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := helper.SendNotificationEmail(verificationRequest.Email, "", "register", "", "", "", false, 1)
		if err != nil {
			fmt.Println("error sending notification email:", err)
		}
	}()

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"OTP verification", nil))
}

func ResetPasswordUser(c echo.Context) error {
	var resetRequest web.ResetRequest
	if err := c.Bind(&resetRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request"))
	}

	if err := helper.ValidateStruct(resetRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Verify OTP
	if err := helper.VerifyOTPByEmail(resetRequest.Email, resetRequest.OTP, "user"); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionGet+"OTP verification failed"))
	}

	hashedPassword := helper.HashPassword(resetRequest.Password)

	// Update password and mark the user as verified
	if err := helper.UpdatePasswordAndMarkVerified(configs.DB, "users", resetRequest.Email, hashedPassword, resetRequest.OTP); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"update password"))
	}

	// Delete OTP from the database
	if err := helper.DeleteOTPFromDatabase(configs.DB, "users", resetRequest.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"delete OTP"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"user's password", nil))
}

func GetOTPForPasswordUser(c echo.Context) error {
	var OTPRequest web.PasswordResetRequest
	if err := c.Bind(&OTPRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(OTPRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	if err := helper.SendOTPViaEmail(OTPRequest.Email, "user", "reset"); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"send OTP"))
	}

	// Use the same success response for both cases
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionCreated+"OTP", nil))
}

func VerifyOTPUser(c echo.Context) error {
	var verificationRequest web.OTPVerificationRequest
	if err := c.Bind(&verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid request"))
	}

	if err := helper.ValidateStruct(verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Verify OTP and handle errors
	if err := helper.VerifyOTPByEmail(verificationRequest.Email, verificationRequest.OTP, "user"); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionGet+"OTP not found"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"OTP verification", nil))
}
