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
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllArticlesByTitleValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/article?title=coba&limit=10&offset=0", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := GetAllArticlesByTitle(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserGetAllArticleValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 5
	url := fmt.Sprintf("/users/articles?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserGetAllArticleMissingLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	url := fmt.Sprintf("/users/articles?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserGetAllArticleMissingOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	limit := 5
	url := fmt.Sprintf("/users/articles?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	fmt.Println(rec.Code)
	err := GetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserGetArticleByIdValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/articles/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:article_id/")
	c.SetParamNames("article_id")
	c.SetParamValues("2")
	err := GetArticleByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUserGetArticleByIdInvalid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/users/articles/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:article_id/")
	c.SetParamNames("article_id")
	c.SetParamValues("99999")
	err := GetArticleByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDoctorGetAllArticlesValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	limit := 5
	url := fmt.Sprintf("/doctors/articles?offset=%d&limit=%d", offset, limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := DoctorGetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDoctorGetAllArticlesMissingOffset(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	limit := 5
	url := fmt.Sprintf("/doctors/articles?limit=%d", limit)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := DoctorGetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDoctorGetAllArticlesMissingLimit(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	url := fmt.Sprintf("/doctors/articles?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := DoctorGetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDoctorGetAllArticlesInvalidId(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	offset := 0
	url := fmt.Sprintf("/doctors/articles?offset=%d", offset)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 9999
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	fmt.Println(rec.Code)
	err := DoctorGetAllArticles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDoctorGetArticleByIdValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/articles/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:article_id/")
	c.SetParamNames("article_id")
	c.SetParamValues("2")
	err := DoctorGetArticleByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDoctorGetArticleByIdInvalid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	req := httptest.NewRequest(http.MethodGet, "/doctors/articles/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	doctorToken := os.Getenv("DOCTOR_TOKEN")
	req.Header.Set("Authorization", doctorToken)
	userID := 1
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)
	c.SetPath("/:article_id/")
	c.SetParamNames("article_id")
	c.SetParamValues("999999999")
	err := DoctorGetArticleByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteArticleByIdValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	articleID := 39
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/path/to/delete/article/%d", articleID), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 15)
	c.SetPath("/:article_id")
	c.SetParamNames("article_id")
	c.SetParamValues(strconv.Itoa(articleID))
	err := DeleteArticleById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCreateArticlesControllerValid(t *testing.T) {
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

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	_ = writer.WriteField("title", "test415")
	_ = writer.WriteField("content", "test415")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/doctors/articles", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 15)
	err = CreateArticle(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestUpdateArticlesControllerValid(t *testing.T) {
	e, db := InitTestDB()
	defer CloseDBTest(db)
	articleID := 18
	imagePath := "../image/gambar.jpg"

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

	_ = writer.WriteField("title", "test223")
	_ = writer.WriteField("content", "test112")

	writer.Close()

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/doctors/articles/%d", articleID), body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", 15)
	c.SetPath("/:article_id")
	c.SetParamNames("article_id")
	c.SetParamValues(strconv.Itoa(articleID))
	err = UpdateArticleById(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
