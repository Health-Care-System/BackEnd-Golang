package controllers

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateComplaintMessageControllerValidImage(t *testing.T) {

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

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")

	err = CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerValidAudio(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)

	userID := 3
	imagePath := "../image/suara.mp3"
	UserToken := os.Getenv("USER_TOKEN")

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")

	err = CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerInvalidUserID(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := "3"
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerUserIDNotSame(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerInvalidRoomchat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 3
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("x")
	err := CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerRoomchatNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 3
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("99999")
	err := CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateComplaintMessageControllerInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 3
	body := `{"message": "success",}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/users/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateComplaintMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}

func TestCreateAdviceMessageControllerValidImage(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)

	doctorID := 4
	imagePath := "../image/gambar.jpg"
	UserToken := os.Getenv("USER_TOKEN")

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")

	err = CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

func TestCreateAdviceMessageControllerValidAudio(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)

	doctorID := 4
	imagePath := "../image/suara.mp3"
	UserToken := os.Getenv("USER_TOKEN")

	file, err := os.Open(imagePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("audio", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	req.Header.Set("Authorization", UserToken)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")

	err = CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
}

func TestCreateAdvicetMessageControllerInvalidDoctorID(t *testing.T) {

	e, db := InitTestDB()
	defer CloseDBTest(db)
	doctorID := "4"
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", doctorID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateAdviceMessageControllerUserIDNotSame(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 5
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateAdviceMessageControllerInvalidRoomchat(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("x")
	err := CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}

func TestCreateAdivceMessageControllerRoomchatNotFound(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	body := `{"message": "success"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("99999")
	err := CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
}

func TestCreateAdviceMessageControllerInvalidBody(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	userID := 4
	body := `{"message": "success",}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/doctors/chats/:roomchat_id/message"), strings.NewReader(body))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:roomchat_id/")
	c.SetParamNames("roomchat_id")
	c.SetParamValues("84")
	err := CreateAdviceMessageController(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
}
