package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"healthcare/models/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLoginAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
        "email":    "ventika20@gmail.com",
        "password": "newpassword123"
    }`

	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestLoginAdminControllerInvalidInput(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestLoginAdminControllerUnregisteredEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{
        "email":    "nonexistent@mail.com",
        "password": "somepassword"
    }`
	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestLoginAdminControllerIncorrectPassword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{
        "email":    "adminsuperr@gmail.com",
        "password": "incorrectpassword"
    }`
	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestLoginAdminControllerEmptyEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{
        "email":    "",
        "password": "somepassword"
    }`
	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestLoginAdminControllerShortPassword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{
        "email":    "adminsuperr@gmail.com",
        "password": "shortpwd"
    }`
	req := httptest.NewRequest(http.MethodPost, "/admins/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginAdminController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAdminProfileControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := 1
	req := httptest.NewRequest(http.MethodGet, "/admins/profile", nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := GetAdminProfileController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetAdminProfileControllerInvalidUserIDType(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidUserID := "invalid"
	req := httptest.NewRequest(http.MethodGet, "/admins/profile", nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", invalidUserID)

	err := GetAdminProfileController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAdminProfileControllerMissingUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAdminProfileController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorsPaymentsByAdminsControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validOffset := 0
	validLimit := 12
	url := fmt.Sprintf("/admins/doctor-payments?offset=%d&limit=%d", validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllDoctorsPaymentsByAdminsController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetAllDoctorsPaymentsByAdminsControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidOffset := -1
	url := fmt.Sprintf("/admins/doctor-payments?offset=%d", invalidOffset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllDoctorsPaymentsByAdminsController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAllDoctorsPaymentsByAdminsControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidLimit := 0
	url := fmt.Sprintf("/admins/doctor-payments?limit=%d", invalidLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllDoctorsPaymentsByAdminsController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorTransactionByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	validTransactionID := "1"

	req := httptest.NewRequest(http.MethodGet, "/admins/doctor-payment", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	q := req.URL.Query()
	q.Add("transaction_id", validTransactionID)
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetDoctorTransactionByIDController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetDoctorTransactionByIDControllerInternalError(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	InvalidTransactionID := "180"

	req := httptest.NewRequest(http.MethodGet, "/admins/doctor-payment", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	q := req.URL.Query()
	q.Add("transaction_id", InvalidTransactionID)
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetDoctorTransactionByIDController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetDoctorTransactionByIDInvalidMissingID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/doctor-transaction", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetDoctorTransactionByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordAdminValid(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "ventika20@gmail.com",
	}

	// Convert struct to JSON string
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = GetOTPForPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetOTPForPasswordAdminInvalidMissingBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/admins/get-otp", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetOTPForPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordAdminInvalidInvalidEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "invalidemail",
	}

	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = GetOTPForPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordAdminInvalidMissingEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = GetOTPForPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordAdminInvalidEmptyEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "",
	}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = GetOTPForPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// func TestVerifyOTPAdminValid(t *testing.T) {

// 	e, db := InitTestDB()
//    defer CloseDBTest(db)
// 	verificationRequest := web.OTPVerificationRequest{
// 		Email: "ventika20@gmail.com",
// 		OTP:   "8895",
// 	}

// 	verificationRequestJSON, err := json.Marshal(verificationRequest)
// 	assert.NoError(t, err)

// 	req := httptest.NewRequest(http.MethodPost, "/admins/verify-otp", bytes.NewReader(verificationRequestJSON))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = VerifyOTPAdmin(c)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// }

func TestVerifyOTPAdminInvalidInvalidOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		Email: "ventika20@gmail.com",
		OTP:   "invalid_otp",
	}

	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = VerifyOTPAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPAdminInvalidMissingEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		// Missing Email field
		OTP: "8895",
	}
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = VerifyOTPAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPAdminInvalidMissingOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		Email: "ventika20@gmail.com",
		// Missing OTP field
	}

	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = VerifyOTPAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// func TestResetPasswordAdminValid(t *testing.T) {

// 	e, db := InitTestDB()
//    defer CloseDBTest(db)
// 	resetRequest := web.ResetRequest{
// 		Email:    "ventika20@gmail.com",
// 		OTP:      "8895",  // Replace with a valid OTP
// 		Password: "newpassword123",
// 	}

// 	// Convert struct to JSON string
// 	resetRequestJSON, err := json.Marshal(resetRequest)
// 	assert.NoError(t, err)

// 	req := httptest.NewRequest(http.MethodPost, "/admins/reset-password", bytes.NewReader(resetRequestJSON))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	err = ResetPasswordAdmin(c)
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusOK, rec.Code)
// }

func TestResetPasswordAdminInvalidInvalidOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	resetRequest := web.ResetRequest{
		Email:    "ventika20@gmail.com",
		OTP:      "invalid_otp",
		Password: "newpassword123",
	}

	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/reset-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordAdminInvalidMissingEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	resetRequest := web.ResetRequest{
		// Missing Email field
		OTP:      "123456",
		Password: "newpassword123",
	}

	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/reset-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordAdminInvalidMissingOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	resetRequest := web.ResetRequest{
		Email: "ventika20@gmail.com",
		// Missing OTP field
		Password: "newpassword123",
	}

	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/reset-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestResetPasswordAdminInvalidInvalidPassword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	resetRequest := web.ResetRequest{
		Email:    "ventika20@gmail.com",
		OTP:      "123456",
		Password: "short",
	}

	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admins/reset-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordAdmin(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	updatedData := web.AdminUpdateRequest{
		Name:     "Updated Admin",
		Email:    "updatedadmin@mail.com",
		Password: "newpassword123",
	}

	req := httptest.NewRequest(http.MethodPut, "/admins/update", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, echo.NewResponse(httptest.NewRecorder(), e))
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updatedData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c.Set("userID", 1)
	err = UpdateAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdatePaymentStatusByAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	updatedData := web.UpdatePaymentRequest{
		PaymentStatus: "pending",
	}

	req := httptest.NewRequest(http.MethodPut, "/admins/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, echo.NewResponse(httptest.NewRecorder(), e))
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	userID := 1
	body, err := json.Marshal(updatedData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c.Set("userID", userID)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("17")
	err = UpdateAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
