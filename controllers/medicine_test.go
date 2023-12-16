package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"healthcare/configs"
	"healthcare/models/web"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func InitTestDB() *echo.Echo {
	e := echo.New()
	godotenv.Load(".env")
	configs.ConnectDBTest()
	return e
}

func TestGetMedicineControllerValid(t *testing.T) {
	e := InitTestDB()
	offset := 0
	limit := 10
	url := fmt.Sprintf("/users/medicines?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineControllerNotFound(t *testing.T) {
	e := InitTestDB()
	offset := 99
	limit := 10
	url := fmt.Sprintf("/users/medicines?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMedicineControllerInvalidOffset(t *testing.T) {
	e := InitTestDB()
	limit := 10
	url := fmt.Sprintf("/users/medicines?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineControllerInvalidLimit(t *testing.T) {
	e := InitTestDB()
	offset := 0
	url := fmt.Sprintf("/users/medicines?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineUserController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineByIDControllerValid(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/users/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("8")
	err := GetMedicineUserByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineByIDControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/users/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err := GetMedicineUserByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineByIDControllerNotFound(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/users/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err := GetMedicineUserByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMedicineAdminControllerValid(t *testing.T) {
	e := InitTestDB()
	offset := 0
	limit := 10
	url := fmt.Sprintf("/admins/medicines?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineAdminControllerNotFound(t *testing.T) {
	e := InitTestDB()
	offset := 99
	limit := 10
	url := fmt.Sprintf("/admins/medicines?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMedicineAdminControllerInvalidOffset(t *testing.T) {
	e := InitTestDB()
	limit := 10
	url := fmt.Sprintf("/admins/medicines?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineAdminControllerInvalidLimit(t *testing.T) {
	e := InitTestDB()
	offset := 0
	url := fmt.Sprintf("/admins/medicines?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetMedicineAdminController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineAdminByIDControllerValid(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("2")
	err := GetMedicineAdminByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetMedicineAdminByIDControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err := GetMedicineAdminByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMedicineAdminByIDControllerNotFound(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err := GetMedicineAdminByIDController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateMedicineControllerBadRequest(t *testing.T) {
	e := InitTestDB()

	// Create a sample MedicineRequest
	medicineRequest := web.MedicineRequest{
		Code: "ABC123",
	}

	_, err := json.Marshal(medicineRequest)
	assert.NoError(t, err)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	url := "/admins/medicines"

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetRequest(req)

	err = CreateMedicineController(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateMedicineAdminControllerValid(t *testing.T) {
	e := InitTestDB()
	updateData := web.MedicineUpdateRequest{
		Name: "Updated Name",
	}
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updateData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("2")
	err = UpdateMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateMedicineAdminControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	updateData := web.MedicineUpdateRequest{
		Name: "Updated Name",
	}
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updateData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err = UpdateMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateMedicineAdminControllerNotFound(t *testing.T) {
	e := InitTestDB()
	updateData := web.MedicineUpdateRequest{
		Name: "Updated Name",
	}
	req := httptest.NewRequest(http.MethodPut, "/admins/medicines", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updateData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err = UpdateMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateImageMedicineAdminControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	updateData := web.MedicineUpdateRequest{
		Name: "Updated Name",
	}
	req := httptest.NewRequest(http.MethodPut, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updateData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err = UpdateImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateImageMedicineAdminControllerNotFound(t *testing.T) {
	e := InitTestDB()
	updateData := web.MedicineUpdateRequest{
		Name: "Updated Name",
	}
	req := httptest.NewRequest(http.MethodPut, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	body, err := json.Marshal(updateData)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err = UpdateImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteImageMedicineAdminByIDControllerInternalServerError(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("22")
	err := DeleteImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDeleteMedicineAdminControllerValid(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("8")
	err := DeleteMedicineController(c)
	if err != nil {
		t.Logf("Error in DeleteMedicineController: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteMedicineAdminControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err := DeleteMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteMedicineAdminControllerNotFound(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err := DeleteMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteImageMedicineAdminControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err := DeleteImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteImageMedicineAdminControllerNotFound(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodDelete, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err := DeleteImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetImageMedicineAdminByIDControllerValid(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("16")
	err := GetImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetImageMedicineAdminByIDControllerInvalidID(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("invalid_id")
	err := GetImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetImageMedicineAdminByIDControllerNotFound(t *testing.T) {
	e := InitTestDB()
	req := httptest.NewRequest(http.MethodGet, "/admins/:medicine_id/medicines/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	AdminToken := os.Getenv("ADMIN_TOKEN")
	req.Header.Set("Authorization", AdminToken)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:medicine_id/")
	c.SetParamNames("medicine_id")
	c.SetParamValues("999")
	err := GetImageMedicineController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
