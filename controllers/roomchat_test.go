package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllDoctorRoomchatController(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("offset", "0")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllDoctorRoomchatControllerWithName(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("offset", "0")
	c.QueryParams().Add("fullname", "Fildza")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetAllDoctorRoomchatControllerWithNameNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("offset", "0")
	c.QueryParams().Add("fullname", "zzzz")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAllDoctorRoomchatControllerNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("offset", "99999")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAllDoctorRoomchatControllerInvalidLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "x")
	c.QueryParams().Add("offset", "0")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAllDoctorRoomchatControllerInvalidOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 4)
	c.QueryParams().Add("limit", "5")
	c.QueryParams().Add("offset", "x")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetAllDoctorRoomchatControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/roomchats", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", "4")
	c.QueryParams().Add("limit", "10")
	c.QueryParams().Add("offset", "0")
	err := GetAllDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUserRoomchatControllerByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", userID)
	err := GetUserRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUserRoomchatControllerByIDControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := "4"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", userID)
	err := GetUserRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUserRoomchatControllerByIDControllerUserIDNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 99999
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", userID)
	err := GetUserRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUserRoomchatControllerByIDControllerInvalidRoomchatID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("x")
	c.Set("userID", userID)
	err := GetUserRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetUserRoomchatControllerByIDControllerRoomchatIDNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("999999")
	c.Set("userID", userID)
	err := GetUserRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetDoctorRoomchatControllerByIDControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", doctorID)
	err := GetDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetDoctorRoomchatControllerByIDControllerInvalidDoctorID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorID := "4"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", doctorID)
	err := GetDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetDoctorRoomchatControllerByIDControllerDoctorIDNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorID := 99999
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("1")
	c.Set("userID", doctorID)
	err := GetDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetDoctorRoomchatControllerByIDControllerInvalidRoomchatID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("x")
	c.Set("userID", doctorID)
	err := GetDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetDoctorRoomchatControllerByIDControllerRoomchatIDNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/chats/:roomchat_id/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorID := 4
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("999999")
	c.Set("userID", doctorID)
	err := GetDoctorRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateRoomchatControllerConflict(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/users/chats/:transaction_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("1")
	c.Set("userID", userID)

	err := CreateRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCreateRoomchatControllerInvalidUserID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/users/chats/:transaction_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := "1"
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("1")
	c.Set("userID", userID)

	err := CreateRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateRoomchatControllerInvalidTransactionID(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/users/chats/:transaction_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("x")
	c.Set("userID", userID)

	err := CreateRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateRoomchatControllerTransactionIDNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/users/chats/:transaction_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("9999")
	c.Set("userID", userID)

	err := CreateRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateRoomchatControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)

	req := httptest.NewRequest(http.MethodPost, "/users/chats/:transaction_id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	userID := 3
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:transaction_id")
	c.SetParamNames("transaction_id")
	c.SetParamValues("153")
	c.Set("userID", userID)

	err := CreateRoomchatController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}
