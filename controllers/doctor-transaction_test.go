package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetDoctorTransactionValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:transaction_id/")
	c.SetParamNames("transaction_id")
	c.SetParamValues("4")
	err := GetDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetDoctorTransactionInvalidInvalidID(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:transaction_id/")
	c.SetParamNames("transaction_id")
	c.SetParamValues("test")
	err := GetDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorTransactionInvalidTransactionID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments/20", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := GetDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorTransactionInternalError(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:transaction_id/")
	c.SetParamNames("transaction_id")
	c.SetParamValues("45")

	err := GetDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorTransactionsControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 3)

	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetAllDoctorTransactionsControllerWithPaymentStatusValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("payment_status", "success")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetAllDoctorTransactionsControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", "invalid_user_id")

	// Call the controller function
	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorTransactionsControllerInvalidPaymentStatus(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("payment_status", "suksesa")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAllDoctorTransactionsControllerNoPaymentStatus(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetAllDoctorTransactionsControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "invalid_limit")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 3)

	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorTransactionsControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "invalid_offset")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorTransactionsControllerInvalid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("payment_status", "non_existent_payment_status")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetAllDoctorTransactionsController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
