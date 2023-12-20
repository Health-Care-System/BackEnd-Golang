package controllers

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatbotValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
		"request": "cara menangani pasien COVID"
	}`
	req := httptest.NewRequest(http.MethodPost, "/chatbot", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := Chatbot(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChatbotInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
		"request": "cara menangani pasien COVID",
	}`
	req := httptest.NewRequest(http.MethodPost, "/chatbot", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := Chatbot(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChatbotRequestRequired(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/chatbot", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := Chatbot(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCustomerService_Payment(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"request": "pembayaran obat"}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCustomerService_RatingDoctor(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"request": "rating dokter"}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCustomerService_ConsultationHistory(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"request": "riwayat konsultasi"}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCustomerService_UnknownRequest(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"request": "other"}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCustomerService_InvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{"request": "other",}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCustomerService_RequestRequired(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{}`
	req := httptest.NewRequest(http.MethodPost, "/customer-service", bytes.NewBufferString(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := CustomerService(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
