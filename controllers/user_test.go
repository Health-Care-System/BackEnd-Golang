package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"healthcare/models/web"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLoginUserControllerValid(t *testing.T) {
	e := InitTestDB()
	requestBody := `{
		"email":    "aditya@gmail.com",
		"password": "aditya12345"
	}`
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginUserControllerInvalid(t *testing.T) {
	e := InitTestDB()
	requestBody := `{
		"email":    "aditya@gmail.com",
		"password": "aditya123456"
    }`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginUserController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginUserControllerInvalidInput(t *testing.T) {
	e := InitTestDB()
	requestBody := `{
		"email":    "",
		"password": "aditya12345"
    }`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginUserController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginUserControllerInvalidEmailFormat(t *testing.T) {
	e := InitTestDB()

	// Membuat email dengan format yang tidak valid
	requestBody := `{
        "email":    "email_tidak_valid",
        "password": "aditya1234"
    }`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginUserController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginUserControllerInvalidPassword(t *testing.T) {
	e := InitTestDB()
	// Membuat password yang tidak valid
	requestBody := `{
        "email":    "aditya@gmail.com",
        "password": "password_tidak_valid"
    }`
	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserControllerValid(t *testing.T) {
	e := InitTestDB()

	// Simulasikan autentikasi dengan menambahkan userID ke konteks
	userID := 1
	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Memanggil controller
	err := GetUserController(c)

	// Memastikan tidak ada kesalahan dan respons memiliki status http.StatusOK
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserControllerInvalidUserIDType(t *testing.T) {
	e := InitTestDB()
	invalidUserID := "invalid"
	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", invalidUserID)

	err := GetUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUserControllerMissingUserID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAllUserByAdminControllerValid(t *testing.T) {
	e := InitTestDB()
	validOffset := 0
	validLimit := 12
	url := fmt.Sprintf("/admins/users?offset=%d&limit=%d", validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllUserByAdminControllerInvalidOffset(t *testing.T) {
	e := InitTestDB()
	invalidOffset := -1
	url := fmt.Sprintf("/admins/users?offset=%d", invalidOffset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := GetAllUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAllUserByAdminControllerInvalidLimit(t *testing.T) {
	e := InitTestDB()
	invalidLimit := 0
	url := fmt.Sprintf("/admins/users?limit=%d", invalidLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Calling the controller
	err := GetAllUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserIDbyAdminControllerValid(t *testing.T) {
	e := InitTestDB()

	req := httptest.NewRequest(http.MethodGet, "/admins/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:user_id/")
	c.SetParamNames("user_id")
	c.SetParamValues("1")
	err := GetUserIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserIDbyAdminControllerMissingIDParam(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/users/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetUserIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserIDbyAdminControllerInvalidIDFormat(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/users/invalid_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("invalid_id")

	err := GetUserIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserIDbyAdminControllerDatabaseError(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/users/999", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("999")

	err := GetUserIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteUserControllerValid(t *testing.T) {
	e := InitTestDB()
	userID := 1
	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Memanggil controller
	err := DeleteUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteUserControllerMissingUserID(t *testing.T) {
	e := InitTestDB()

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := DeleteUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code) // Assuming a 500 status for missing userID
}

func TestDeleteUserByAdminControllerValid(t *testing.T) {
	e := InitTestDB()

	req := httptest.NewRequest(http.MethodDelete, "/admins/user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:user_id/")
	c.SetParamNames("user_id")
	c.SetParamValues("1")
	err := DeleteUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteUserByAdminControllerMissingIDParam(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/users/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := DeleteUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteUserByAdminControllerInvalidIDFormat(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/users/invalid_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("invalid_id")

	err := DeleteUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteUserByAdminControllerDatabaseError(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/users/999", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("user_id")
	c.SetParamValues("999")

	err := DeleteUserByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetOTPForPasswordUserValid(t *testing.T) {
	// Initialize Echo and create a fake context
	e := InitTestDB()
	otpRequest := web.PasswordResetRequest{
		Email: "patient1test@gmail.com",
	}

	// Convert struct to JSON string
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetOTPForPasswordUserInvalidMissingBody(t *testing.T) {
	e := InitTestDB()

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOTPForPasswordUserInvalidInvalidEmail(t *testing.T) {

	e := InitTestDB()
	otpRequest := web.PasswordResetRequest{
		Email: "invalidemail",
	}

	// Convert struct to JSON string
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOTPForPasswordUserInvalidMissingEmail(t *testing.T) {

	e := InitTestDB()
	otpRequest := web.PasswordResetRequest{}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOTPForPasswordUserInvalidEmptyEmail(t *testing.T) {

	e := InitTestDB()
	otpRequest := web.PasswordResetRequest{
		Email: "",
	}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOTPForPasswordUserInvalidEmptyRequestBody(t *testing.T) {
	e := InitTestDB()

	req := httptest.NewRequest(http.MethodPost, "/users/get-otp", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetOTPForPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPUserValid(t *testing.T) {

	e := InitTestDB()
	verificationRequest := web.OTPVerificationRequest{
		Email: "patient1test@gmail.com",
		OTP:   "4159",
	}

	// Convert struct to JSON string
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestVerifyOTPUserInvalidInvalidOTP(t *testing.T) {

	e := InitTestDB()
	verificationRequest := web.OTPVerificationRequest{
		Email: "patient1test@gmail.com",
		OTP:   "invalid_otp",
	}

	// Convert struct to JSON string
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPUserInvalidMissingEmail(t *testing.T) {
	e := InitTestDB()
	verificationRequest := web.OTPVerificationRequest{
		// Missing Email field
		OTP: "4159",
	}
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPUserInvalidMissingOTP(t *testing.T) {
	e := InitTestDB()
	verificationRequest := web.OTPVerificationRequest{
		Email: "patient1test@gmail.com",
		// Missing OTP field
	}

	// Convert struct to JSON string
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPUserInvalidEmptyEmail(t *testing.T) {

	e := InitTestDB()
	verificationRequest := web.OTPVerificationRequest{
		Email: "", // Empty Email field
		OTP:   "4159",
	}
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordUserValid(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		Email:    "patient1test@gmail.com",
		OTP:      "4159",
		Password: "userpass123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestResetPasswordUserInvalidInvalidOTP(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		Email:    "patient1test@gmail.com",
		OTP:      "invalid_otp",
		Password: "userpass123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordUserInvalidMissingEmail(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		// Missing Email field
		OTP:      "4159",
		Password: "userpass123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordUserInvalidMissingOTP(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		Email: "patient1test@gmail.com",
		// Missing OTP field
		Password: "userpass123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordUserInvalidInvalidPassword(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		Email:    "patient1test@gmail.com",
		OTP:      "4159",
		Password: "short",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordUserInvalidEmptyEmail(t *testing.T) {
	e := InitTestDB()

	resetRequest := web.ResetRequest{
		Email:    "", // Empty Email field
		OTP:      "4159",
		Password: "userpass123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/users/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

