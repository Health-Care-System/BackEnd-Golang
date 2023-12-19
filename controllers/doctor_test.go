package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"healthcare/models/web"
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

func TestRegisterDoctorControllerValidWithImage(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	imagePath := "../image/gambar.jpg"

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Add image file to the form data
	part, err := writer.CreateFormFile("profile_picture", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	_ = writer.WriteField("fullname", "Hanisah Fildza Annafisah")
	_ = writer.WriteField("email", "asal3@gmail.com")
	_ = writer.WriteField("password", "asalasalan")
	_ = writer.WriteField("price", "70000")
	_ = writer.WriteField("no_str", "12345678910")
	_ = writer.WriteField("specialist", "dokter hewan")
	_ = writer.WriteField("gender", "female")
	_ = writer.WriteField("experience", "5 tahun")
	_ = writer.WriteField("alumnus", "UI")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admins/doctors/register", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = RegisterDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestRegisterDoctorControllerInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
            "fullname": "Hanisah Fildza Annafisah",
            "email": "asal@gmail.com",
            "password": "asalasalan",
            "price": "70000",
            "no_str": 12345678910,
			"specialist": "dokter hewan",
			"gender": "male",
			"exprience": "5 tahun",
			"alumnus": "UI",
	}`
	req := httptest.NewRequest(http.MethodPost, "/admins/doctors/register", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := RegisterDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterDoctorControllerRequiredFields(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
           "fullname": "Hanisah Fildza Annafisah",
           "email": "asal2@gmail.com",
           "password": "asalasalan"
	}`
	req := httptest.NewRequest(http.MethodPost, "/admins/doctors/register", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := RegisterDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateDoctorControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	imagePath := "../image/gambar.jpg"
	doctorID := 9

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("profile_picture", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	_ = writer.WriteField("email", "asalasalid9@gmail.com")
	_ = writer.WriteField("password", "asalasalan")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/doctors/profile/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err = UpdateDoctorController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateDoctorByAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	imagePath := "../image/gambar.jpg"

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("profile_picture", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}
	_ = writer.WriteField("email", "asalasalid10@gmail.com")
	_ = writer.WriteField("password", "asalasalan")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/admins/doctor/:doctor_id/", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("10")
	err = UpdateDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateManageUserControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 1
	requestBody := `{"health_details": "demam",
					 "patient_status": "recovered"}`
	req := httptest.NewRequest(http.MethodPut, "/doctors/manage-user/:transaction_id/", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("1")
	err := UpdateManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginDoctorControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
		"email":    "mutiakhoirunniza@gmail.com",
		"password": "asalasalan"
	}`
	req := httptest.NewRequest(http.MethodPost, "/doctors/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestLoginDoctorControllerInvalid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	requestBody := `{
        "email":    "mutiakhoirunniza@gmail.com",
        "password": "cokayaa1234"
    }`
	req := httptest.NewRequest(http.MethodPost, "/doctor/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestLoginDoctorControllerInvalidInput(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	requestBody := `{
        "email":    "",
        "password": "cokayaa123"
    }`
	req := httptest.NewRequest(http.MethodPost, "/doctor/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestLoginDoctorControllerInvalidEmailFormat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	// Membuat email dengan format yang tidak valid
	requestBody := `{
        "email":    "email_tidak_valid",
        "password": "cokayaa123"
    }`
	req := httptest.NewRequest(http.MethodPost, "/doctor/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestLoginDoctorControllerInvalidPassword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	// Membuat password yang tidak valid
	requestBody := `{
        "email":    "mutiakhoirunniza@gmail.com",
        "password": "password_tidak_valid"
    }`
	req := httptest.NewRequest(http.MethodPost, "/doctor/login", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := LoginDoctorController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAvailableDoctorControllervalid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 12
	url := fmt.Sprintf("/users/doctors/available?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetAvailableDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllDoctorConsultationControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 4
	offset := 0
	limit := 12
	url := fmt.Sprintf("/doctors/doctor-consultation?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetAllDoctorConsultationController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllDoctorConsultationControllerInvalidDoctorID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := "4"
	offset := 0
	limit := 12
	url := fmt.Sprintf("/doctors/doctor-consultation?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetAllDoctorConsultationController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAllDoctorConsultationControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 4
	limit := 12
	url := fmt.Sprintf("/doctors/doctor-consultation?&limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetAllDoctorConsultationController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAllDoctorConsultationControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 4
	offset := 0
	url := fmt.Sprintf("/doctors/doctor-consultation?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetAllDoctorConsultationController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetDoctorStatusControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 4
	url := fmt.Sprintf("/doctors/status")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetDoctorStatusControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 99999
	url := fmt.Sprintf("/doctors/status")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetDoctorStatusControllerInternalServerError(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := 99999
	url := fmt.Sprintf("/doctors/status")
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	err := GetDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAvailableDoctorControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidOffset := -1
	url := fmt.Sprintf("/users/doctors/available?offset=%d", invalidOffset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAvailableDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAvailableDoctorControllerInvalidLimit(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidLimit := -10

	url := fmt.Sprintf("/users/doctors/available?limit=%d", invalidLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAvailableDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetSpecializeDoctorControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validSpecialist := "Bedah"
	validOffset := 0
	validLimit := 12
	url := fmt.Sprintf("/users/doctors?specialist=%s&offset=%d&limit=%d", validSpecialist, validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetSpecializeDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetSpecializeDoctorControllerMissingSpecialist(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validOffset := 0
	validLimit := 12
	url := fmt.Sprintf("/users/doctors?offset=%d&limit=%d", validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetSpecializeDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetSpecializeDoctorControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validSpecialist := "Bedah"
	invalidLimit := -5
	url := fmt.Sprintf("/users/doctors?specialist=%s&limit=%d", validSpecialist, invalidLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := GetSpecializeDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetSpecializeDoctorControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validSpecialist := "Bedah"
	validOffset := -10
	url := fmt.Sprintf("/users/doctors?specialist=%s&offset=%d", validSpecialist, validOffset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := GetSpecializeDoctor(c)

	// Memastikan tidak ada kesalahan dan respons memiliki status http.StatusBadRequest
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorProfileControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	// Simulasikan autentikasi dengan menambahkan userID ke konteks
	userID := 1
	req := httptest.NewRequest(http.MethodGet, "/doctors/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Memanggil controller
	err := GetDoctorProfileController(c)

	// Memastikan tidak ada kesalahan dan respons memiliki status http.StatusOK
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetDoctorProfileControllerInvalidUserIDType(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidUserID := "invalid"
	req := httptest.NewRequest(http.MethodGet, "/doctors/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", invalidUserID)

	err := GetDoctorProfileController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetDoctorProfileControllerMissingUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/profile", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetDoctorProfileController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetAllDoctorByAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validOffset := 0
	validLimit := 12
	url := fmt.Sprintf("/admins/doctors?offset=%d&limit=%d", validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllDoctorByAdminControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	validOffset := 999999
	validLimit := 12
	url := fmt.Sprintf("/admins/doctors?offset=%d&limit=%d", validOffset, validLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAllDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAllDoctorByAdminControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidOffset := -1
	url := fmt.Sprintf("/admins/doctors?offset=%d", invalidOffset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := GetAllDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetAllDoctorByAdminControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	invalidLimit := 0
	url := fmt.Sprintf("/admins/doctors?limit=%d", invalidLimit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Calling the controller
	err := GetAllDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestDeleteDoctorControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 1
	req := httptest.NewRequest(http.MethodDelete, "/doctors", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Memanggil controller
	err := DeleteDoctorController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestDeleteDoctorControllerMissingUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodDelete, "/doctors", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Memanggil controller
	err := DeleteDoctorController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code) // Assuming a 500 status for missing userID
}
func TestGetDoctorByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("4")
	fmt.Println(rec.Code)

	err := GetDoctorByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetDoctorByIDControllerInvalidIDParam(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("invalid_id")
	fmt.Println(rec.Code)

	err := GetDoctorByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorByIDControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("999")
	fmt.Println(rec.Code)

	err := GetDoctorByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetDoctorByIDControllerInvalidIDFormat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("abc") // Non-numeric ID
	fmt.Println(rec.Code)

	err := GetDoctorByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorByIDControllerMissingIDParam(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Omitting the doctor_id parameter
	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	// Not setting ParamValues intentionally to simulate missing parameter
	fmt.Println(rec.Code)

	err := GetDoctorByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestDeleteDoctorByAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodDelete, "/admins/doctor", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("16")
	err := DeleteDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestDeleteDoctorByAdminControllerMissingIDParam(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/admins/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := DeleteDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestDeleteDoctorByAdminControllerInvalidIDFormat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/admins/doctors/invalid_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("doctor_id")
	c.SetParamValues("invalid_id")

	err := DeleteDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestDeleteDoctorByAdminControllerDatabaseError(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodDelete, "/admins/doctors/999", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("doctor_id")
	c.SetParamValues("999")

	err := DeleteDoctorByAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetDoctorIDbyAdminControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/admins/doctor", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:doctor_id/")
	c.SetParamNames("doctor_id")
	c.SetParamValues("4")
	err := GetDoctorIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetDoctorIDbyAdminControllerMissingIDParam(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/doctors/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetDoctorIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorIDbyAdminControllerInvalidIDFormat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/doctors/invalid_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("doctor_id")
	c.SetParamValues("invalid_id")

	err := GetDoctorIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetDoctorIDbyAdminControllerDatabaseError(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/admins/doctors/999", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("doctor_id")
	c.SetParamValues("999")

	err := GetDoctorIDbyAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
func TestGetManageUserControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetManageUserControllerInvalidLimit(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "invalid_limit")
	q.Add("offset", "0")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetManageUserControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	// Set up a request with an invalid offset parameter
	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "invalid_offset")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetManageUserControllerNoResults(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("transaction_id", "999")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetManageUserControllerNonExistentFullname(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("fullname", "non_existent_fullname")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	// Call the controller function
	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetManageUserControllerInvalidPatientStatus(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("patient_status", "invalid_patient_status")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("userID", 1)

	err := GetManageUserController(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetManageUserControllerInvalidFullname(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("fullname", "invalid_fullname")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetManageUserControllerInvalidKeyword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodGet, "/doctors/manage-user", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)

	q := req.URL.Query()
	q.Add("limit", "10")
	q.Add("offset", "0")
	q.Add("keyword", "invalid_keyword")
	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 1)

	err := GetManageUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
func TestGetOTPForPasswordDoctorValid(t *testing.T) {
	// Initialize Echo and create a fake context
	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "mutiakhoirunniza@gmail.com",
	}

	// Convert struct to JSON string
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
func TestGetOTPForPasswordDoctorInvalidMissingBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordDoctorInvalidInvalidEmail(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "invalidemail",
	}

	// Convert struct to JSON string
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordDoctorInvalidMissingEmail(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordDoctorInvalidEmptyEmail(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	otpRequest := web.PasswordResetRequest{
		Email: "",
	}
	otpRequestJSON, err := json.Marshal(otpRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", bytes.NewReader(otpRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestGetOTPForPasswordDoctorInvalidEmptyRequestBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/doctors/get-otp", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetOTPForPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

//func TestVerifyOTPDoctorValid(t *testing.T) {
//
//	e, db := InitTestDB()
//    defer CloseDBTest(db)
//	verificationRequest := web.OTPVerificationRequest{
//		Email: "mutiakhoirunniza@gmail.com",
//		OTP:   "5023",
//	}
//
//	// Convert struct to JSON string
//	verificationRequestJSON, err := json.Marshal(verificationRequest)
//	assert.NoError(t, err)
//
//	req := httptest.NewRequest(http.MethodPost, "/doctors/verify-otp", bytes.NewReader(verificationRequestJSON))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//	err = VerifyOTPDoctor(c)
//	assert.NoError(t, err)
//	assert.Equal(t, http.StatusOK, rec.Code)
//}

func TestVerifyOTPDoctorInvalidInvalidOTP(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		Email: "mutiakhoirunniza@gmail.com",
		OTP:   "invalid_otp",
	}

	// Convert struct to JSON string
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestVerifyOTPDoctorInvalidMissingEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		// Missing Email field
		OTP: "5023",
	}
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestVerifyOTPDoctorInvalidMissingOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		Email: "mutiakhoirunniza@gmail.com",
		// Missing OTP field
	}

	// Convert struct to JSON string
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestVerifyOTPDoctorInvalidEmptyEmail(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	verificationRequest := web.OTPVerificationRequest{
		Email: "", // Empty Email field
		OTP:   "5023",
	}
	verificationRequestJSON, err := json.Marshal(verificationRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/verify-otp", bytes.NewReader(verificationRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = VerifyOTPDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

//		func TestResetPasswordDoctorValid(t *testing.T) {
//			e, db := InitTestDB()
//	   defer CloseDBTest(db)
//
//			resetRequest := web.ResetRequest{
//				Email:    "mutiakhoirunniza@gmail.com",
//				OTP:      "5023",
//				Password: "newpassword123",
//			}
//
//			// Convert struct to JSON string
//			resetRequestJSON, err := json.Marshal(resetRequest)
//			assert.NoError(t, err)
//
//			req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
//			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//			rec := httptest.NewRecorder()
//			c := e.NewContext(req, rec)
//
//			err = ResetPasswordDoctor(c)
//			assert.NoError(t, err)
//			assert.Equal(t, http.StatusOK, rec.Code)
//		}
func TestResetPasswordDoctorInvalidInvalidOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	resetRequest := web.ResetRequest{
		Email:    "mutiakhoirunniza@gmail.com",
		OTP:      "invalid_otp",
		Password: "newpassword123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestResetPasswordDoctorInvalidMissingEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	resetRequest := web.ResetRequest{
		// Missing Email field
		OTP:      "5023",
		Password: "newpassword123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestResetPasswordDoctorInvalidMissingOTP(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	resetRequest := web.ResetRequest{
		Email: "mutiakhoirunniza@gmail.com",
		// Missing OTP field
		Password: "newpassword123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestResetPasswordDoctorInvalidInvalidPassword(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	resetRequest := web.ResetRequest{
		Email:    "mutiakhoirunniza@gmail.com",
		OTP:      "5023",
		Password: "short",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
func TestResetPasswordDoctorInvalidEmptyEmail(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	resetRequest := web.ResetRequest{
		Email:    "", // Empty Email field
		OTP:      "5023",
		Password: "newpassword123",
	}

	// Convert struct to JSON string
	resetRequestJSON, err := json.Marshal(resetRequest)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/doctors/change-password", bytes.NewReader(resetRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = ResetPasswordDoctor(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChangeDoctorStatusControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorToken := os.Getenv("DOCTOR_TOKEN")

	requestBody := `{"status": true}`
	req := httptest.NewRequest(http.MethodPut, "/doctors/status", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", doctorToken)
	doctorID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)

	err := ChangeDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestChangeDoctorStatusControllerInvalidDoctorID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorToken := os.Getenv("DOCTOR_TOKEN")

	requestBody := `{"status": true}`
	req := httptest.NewRequest(http.MethodPut, "/doctors/status", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", doctorToken)
	doctorID := "4"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)

	err := ChangeDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestChangeDoctorStatusControllerInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorToken := os.Getenv("DOCTOR_TOKEN")

	requestBody := `{"status": true,}`
	req := httptest.NewRequest(http.MethodPut, "/doctors/status", strings.NewReader(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", doctorToken)
	doctorID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)

	err := ChangeDoctorStatusController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
