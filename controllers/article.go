package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/helper/constanta"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateArticle(c echo.Context) error {
	userID := c.Get("userID").(int)

	var article web.CreateArticle

	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}
	if err := helper.ValidateStruct(article); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrImageFileRequired))
	}
	defer file.Close()

	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	ext := filepath.Ext(fileHeader.Filename)
	allowed := false
	for _, validExt := range allowedExtensions {
		if ext == validExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidImageFormat))
	}

	image, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to Cloud Storage"))
	}

	articleRequest := request.ConvertToCreateArticleRequest(article)

	articleRequest.Image = image
	articleRequest.DoctorID = uint(userID)

	if err := configs.DB.Create(&articleRequest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to create article"))
	}

	return c.JSON(http.StatusCreated, helper.SuccessResponse(constanta.SuccessActionCreated+"article", nil))
}

func UpdateArticleById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	userID := c.Get("userID").(int)

	var existingArticle schema.Article
	result := configs.DB.First(&existingArticle, id, userID)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	var updateArticle web.CreateArticle

	if err := c.Bind(&updateArticle); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}
	if err := helper.ValidateStruct(updateArticle); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	if file, fileHeader, err := c.Request().FormFile("image"); err == nil {
		defer file.Close()

		allowedExtensions := []string{".jpg", ".jpeg", ".png"}
		ext := filepath.Ext(fileHeader.Filename)
		allowed := false
		for _, validExt := range allowedExtensions {
			if ext == validExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidImageFormat))
		}

		image, err := helper.UploadFilesToGCS(c, fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to Cloud Storage"))
		}

		updateArticle.Image = image

	} else if err != http.ErrMissingFile {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrImageFileRequired))

	} else {
		updateArticle.Image = existingArticle.Image
	}

	result = configs.DB.Model(&existingArticle).Updates(updateArticle)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"article"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"article", nil))

}

func GetAllArticles(c echo.Context) error {
	var articles []schema.Article

	err := configs.DB.Find(&articles).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	if len(articles) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ListConvertToArticleDoctors(articles)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"article", response))
}

func GetArticleByID(c echo.Context) error {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var article schema.Article

	if err = configs.DB.First(&article, articleID).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToArticleDoctor(article)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet, response))
}

func GetAllArticlesByTitle(c echo.Context) error {

	param := c.QueryParam("title")
	if param == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("title"+constanta.ErrQueryParamRequired))
	}

	var articles []schema.Article
	err := configs.DB.Where("title LIKE ?", "%"+param+"%").Find(&articles).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	if len(articles) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ListConvertToArticleDoctors(articles)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"article", response))
}

func DeleteArticleById(c echo.Context) error {
	userID := c.Get("userID").(int)

	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var existingArticle schema.Article
	result := configs.DB.First(&existingArticle, articleID, userID)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	if err := configs.DB.Delete(&existingArticle).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"article"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionDeleted+"article", nil))
}

// get all for doctor
// get by id for doctor

func DoctorGetAllArticles(c echo.Context) error {
	userID := c.Get("userID").(int)

	var articles []schema.Article

	err := configs.DB.Where("doctor_id = ?", userID).Find(&articles).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"articles"))
	}

	if len(articles) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToGetAllArticles(articles)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"articles", response))
}

func DoctorGetArticleByID(c echo.Context) error {
	userID := c.Get("userID").(int)

	articleID := c.Param("id")

	var article schema.Article
	if err := configs.DB.Where("id = ? AND doctor_id = ?", articleID, userID).Find(&article).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToGetArticleResponse(&article)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"article", response))
}
