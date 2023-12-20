package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

func TestGetDoctorTransactionControllerInvalidInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := "3"
	UserToken := os.Getenv("USER_TOKEN")
	req := httptest.NewRequest(http.MethodGet, "/users/doctor-payments/", nil)
	req.Header.Set("Authorization", UserToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
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

func TestCreateDoctorTransactionControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := "invalid_user_id"
	doctorID := 1

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/doctor-payments/%d", doctorID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := CreateDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestCreateDoctorTransactionControllerInvalidDoctorID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := 1
	doctorID := "10"

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/doctor-payments/%s", doctorID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := CreateDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestCreateDoctorTransactionControllerInvalidInputData(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := 1
	doctorID := 1
	requestBody := `{"payment_method": "invalid_method", "payment_confirmation": ""}`

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/doctor-payments/%d", doctorID), strings.NewReader(requestBody))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	err := CreateDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestCreateDoctorTransactionControllerValid(t *testing.T) {
	
	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := 3
	imagePath := "../image/gambar.jpg"
	UserToken := os.Getenv("USER_TOKEN")

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("payment_confirmation", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	
	writer.WriteField("payment_method", "manual transfer bca")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/doctor-payments/:doctor_id/4"), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("4")


	err = CreateDoctorTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

