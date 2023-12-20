package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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
	c.SetParamValues("1")
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
	c.SetParamValues("1")
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
