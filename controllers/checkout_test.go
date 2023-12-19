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

func TestGetUserCheckoutControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments/checkout?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserCheckoutControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments/checkout?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserCheckoutControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 10
	url := fmt.Sprintf("/users/medicines-payments/checkout?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserCheckoutControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments/checkout?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := "test"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserCheckoutControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments/checkout?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 112
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetUserCheckoutControllerNotFoundStatus(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	payment_status := "test"
	url := fmt.Sprintf("/users/medicines-payments/checkout?offset=%d&limit=%d&payment_status=%s", offset, limit, payment_status)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetUserCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetUserCheckoutByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("17")
	err := GetUserCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserCheckoutByIDControllerInvalidID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("test")
	err := GetUserCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserCheckoutByIDControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := "1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("17")
	err := GetUserCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserCheckoutByIDControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 112
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("1")
	err := GetUserCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAdminCheckoutControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	url := fmt.Sprintf("/admins/medicines-payments/checkout?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAdminCheckoutControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	limit := 10
	url := fmt.Sprintf("/admins/medicines-payments/checkout?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAdminCheckoutControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	url := fmt.Sprintf("/admins/medicines-payments/checkout?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAdminCheckoutControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	payment_status := "test"
	url := fmt.Sprintf("/admins/medicines-payments/checkout?offset=%d&limit=%d&payment_status=%s", offset, limit, payment_status)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAdminCheckoutControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	user_id := "test"
	url := fmt.Sprintf("/admins/medicines-payments/checkout?offset=%d&limit=%d&user_id=%s", offset, limit, user_id)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAdminCheckoutControllerNotFoundByUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	user_id := 112
	url := fmt.Sprintf("/admins/medicines-payments/checkout?offset=%d&limit=%d&user_id=%d", offset, limit, user_id)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAdminCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAdminCheckoutByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("1")
	err := GetAdminCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAdminCheckoutByIDControllerInvalidID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("test")
	err := GetAdminCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAdminCheckoutByIDControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines-payments/checkout/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("999999999")
	err := GetAdminCheckoutByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCheckoutControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"payment_status": "success"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines-payments/checkout/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("9")
	err := UpdateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateCheckoutControllerInvalidID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"payment_status": "success"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines-payments/checkout/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("tes")
	err := UpdateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCheckoutControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"payment_status": "success"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines-payments/checkout/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("9999999999")
	err := UpdateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCheckoutControllerInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"payment_status": "success",}`
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines-payments/checkout/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("1")
	err := UpdateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCheckoutControllerInvalidStock(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"payment_status": "success"}`
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines-payments/checkout/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:checkout_id/")
	c.SetParamNames("checkout_id")
	c.SetParamValues("1")
	err := UpdateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCheckoutControllerInvalid(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	medicineTransactionID := 18
	UserToken := os.Getenv("USER_TOKEN")
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/medicines-payments/checkout?medicine_transaction_id=%d", medicineTransactionID), nil)
	req.Header.Set("Authorization", UserToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	err := CreateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCheckoutControllerInvalidUserID(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := "4"
	medicineTransactionID := 18
	UserToken := os.Getenv("USER_TOKEN")
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/medicines-payments/checkout?medicine_transaction_id=%d", medicineTransactionID), nil)
	req.Header.Set("Authorization", UserToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	err := CreateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCheckoutControllerInvalidBind(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	medicineTransactionID := 18
	requestBody := `{"name": "nathan",}]}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/medicines-payments/checkout?medicine_transaction_id=%d", medicineTransactionID), strings.NewReader(requestBody))
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	err := CreateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCheckoutControllerValid(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 1
	medicineTransactionID := 205
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
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/medicines-payments/checkout?medicine_transaction_id=%d", medicineTransactionID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	err = CreateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateCheckoutControllerForbiddenUserID(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	medicineTransactionID := 203
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
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/medicines-payments/checkout?medicine_transaction_id=%d", medicineTransactionID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	err = CreateCheckoutController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}
