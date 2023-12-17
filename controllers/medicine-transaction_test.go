package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetMedicineTransactionControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineTransactionControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineTransactionControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	url := fmt.Sprintf("/users/medicines-payments?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineTransactionControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 10
	status_transaction := "belum dibayar"
	status := url.QueryEscape(status_transaction)
	url := fmt.Sprintf("/users/medicines-payments?offset=%d&limit=%d&status_transaction=%s", offset, limit, status)
	fmt.Println("URL:", url)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineTransactionControllerNotFoundOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 999
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments?offset=%d&limit=%d", offset, limit)
	fmt.Println("URL:", url)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMedicineTransactionControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 999
	limit := 10
	url := fmt.Sprintf("/users/medicines-payments?offset=%d&limit=%d", offset, limit)
	fmt.Println("URL:", url)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := "invalid"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := GetMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineTransactionByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("9")
	c.Set("userID", userID)
	err := GetMedicineTransactionByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineTransactionByIDControllerInvalidID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 41
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("invalid_id")
	c.Set("userID", userID)
	err := GetMedicineTransactionByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineTransactionByIDControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := "test"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("16")
	c.Set("userID", userID)
	err := GetMedicineTransactionByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineTransactionByIDControllerInvalidNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("16")
	c.Set("userID", userID)
	err := GetMedicineTransactionByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteMedicineTransactionControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("10")
	c.Set("userID", userID)
	err := DeleteMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteMedicineTransactionControllerDenied(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 2
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("17")
	c.Set("userID", userID)
	err := DeleteMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteMedicineTransactionControllerInvalidID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/users/medicines-payments/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medtrans_id/")
	c.SetParamNames("medtrans_id")
	c.SetParamValues("invalid_id")
	c.Set("userID", userID)
	err := DeleteMedicineTransactionController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMedicineTransactionControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca",
            "medicine_details": [{
                    "medicine_id": 2,
                    "quantity": 1
                }]}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestCreateMedicineTransactionControllerRequiredFields(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca"}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMedicineTransactionControllerBind(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah"
            "address": "cengkareng"}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMedicineTransactionControllerMedicineIDNil(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca",
            "medicine_details": [{
                    "medicine_id": 0,
                    "quantity": 1
                }]}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMedicineTransactionControllerValidQuantityNil(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca",
            "medicine_details": [{
                    "medicine_id": 1,
                    "quantity": 0
                }]}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateMedicineTransactionControllerMedicineNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "nathan",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca",
            "medicine_details": [{
                    "medicine_id": 999,
                    "quantity": 1
                }]}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateMedicineTransactionControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "name": "Hanisah Fildza Annafisah",
            "address": "cengkareng",
            "hp": "08123",
            "payment_method": "manual transfer bca",
            "medicine_details": [{
                    "medicine_id": 999,
                    "quantity": 1
                }]}`
	req := httptest.NewRequest(http.MethodPost, "/users/medicines-payments", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	UserToken := os.Getenv("USER_TOKEN")
	req.Header.Set("Authorization", UserToken)
	userID := "test"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := CreateMedicineTransaction(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
