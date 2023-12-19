package controllers

import (
	"fmt"
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/helper/constanta"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// Create Medicine
func CreateMedicineController(c echo.Context) error {

	var medicine web.MedicineRequest

	if err := c.Bind(&medicine); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(medicine); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Upload files
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

	imageURL, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to cloud storage"))
	}

	medicine.Image = imageURL

	medicineRequest := request.ConvertToMedicineRequest(medicine)

	if err := configs.DB.Create(&medicineRequest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionCreated+"medicine"))
	}

	response := response.ConvertToAdminMedicineResponse(medicineRequest)

	return c.JSON(http.StatusCreated, helper.SuccessResponse(constanta.SuccessActionCreated+"medicine", response))
}

// Update Medicine by ID
func UpdateMedicineController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var existingMedicine schema.Medicine

	result := configs.DB.First(&existingMedicine, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	var updatedMedicineRequest web.MedicineUpdateRequest

	if err := c.Bind(&updatedMedicineRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(updatedMedicineRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	result = configs.DB.Model(&existingMedicine).Updates(updatedMedicineRequest)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"medicine"))
	}

	response := response.ConvertToAdminMedicineUpdateResponse(&existingMedicine)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"medicine", response))
}

// Update Image Medicine by ID
func UpdateImageMedicineController(c echo.Context) error {
	medicineID, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var existingMedicine schema.Medicine
	if err := configs.DB.First(&existingMedicine, medicineID).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
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

	newImage, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to upload image to cloud storage"))
	}

	if existingMedicine.Image != "" {
		oldFilename := path.Base(existingMedicine.Image)

		if err := helper.DeleteFilesFromGCS(oldFilename); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to delete old image"))
		}
	}

	existingMedicine.Image = newImage

	if err := configs.DB.Save(&existingMedicine).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"image medicine"))
	}

	response := response.ConvertToAdminMedicineImageResponse(&existingMedicine)
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"image medicine", response))
}

// Delete Medicine by ID
func DeleteMedicineController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var medicine schema.Medicine

	result := configs.DB.First(&medicine, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	result = configs.DB.Delete(&medicine, id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"medicine"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionDeleted+"medicine", nil))
}

// Delete Image Medicine by ID
func DeleteImageMedicineController(c echo.Context) error {
	medicineID, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var medicine schema.Medicine
	if err := configs.DB.First(&medicine, medicineID).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	if medicine.Image != "" {
		filename := path.Base(medicine.Image)
		log.Printf("Deleting file from GCS: %s", filename)
		if err := helper.DeleteFilesFromGCS(filename); err != nil {
			log.Printf("Error deleting file from GCS: %v", err)
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"image medicine"))
		}
	}

	medicine.Image = ""

	if err := configs.DB.Save(&medicine).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"image medicine"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionDeleted+"image medicine", nil))
}

// Get Image Medicine by ID
func GetImageMedicineController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var medicine schema.Medicine
	if err := configs.DB.First(&medicine, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToAdminMedicineImageResponse(&medicine)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"image medicine", response))
}

func GetAll(offset, limit int, price, category, keyword string, queryInput []schema.Medicine) ([]schema.Medicine, int64, error) {

	if offset < 0 || limit < 0 {
		return nil, 0, nil
	}

	queryAll := queryInput
	var total int64

	query := configs.DB.Model(&queryAll)

	if keyword != "" {
		query = query.Where("name LIKE ? OR merk LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if price == "low" {
		query = query.Order("price ASC")
	}

	if price == "high" {
		query = query.Order("price DESC")
	}

	if category != "" {
		query = query.Where("category LIKE ?", "%"+category+"%")
	}

	query.Find(&queryAll).Count(&total)

	query = query.Limit(limit).Offset(offset)

	result := query.Find(&queryAll)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	if offset >= int(total) {
		return nil, 0, fmt.Errorf("not found")
	}

	return queryAll, total, nil
}

// Admin Get All Medicines Pagination
func GetMedicineAdminController(c echo.Context) error {
	params := c.QueryParams()
	limit, err := strconv.Atoi(params.Get("limit"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(params.Get("offset"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	keyword := params.Get("keyword")
	price := params.Get("price")
	category := params.Get("category")

	var medicines []schema.Medicine

	medicine, total, err := GetAll(offset, limit, price, category, keyword, medicines)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("medicines "+constanta.ErrNotFound))
		}
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
	}

	pagination := helper.Pagination(offset, limit, total)

	response := response.ConvertToAdminGetAllMedicinesResponse(medicine)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"medicines", response, pagination))
}

// Admin Get Medicine by ID
func GetMedicineAdminByIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var medicine schema.Medicine

	if err := configs.DB.First(&medicine, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToAdminMedicineResponse(&medicine)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"medicine", response))
}

// User Get All Medicines Pagination
func GetMedicineUserController(c echo.Context) error {
	params := c.QueryParams()
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	keyword := params.Get("keyword")
	price := params.Get("price")
	category := params.Get("category")

	var medicines []schema.Medicine

	medicine, total, err := GetAll(offset, limit, price, category, keyword, medicines)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("medicines "+constanta.ErrNotFound))
		}
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
	}

	pagination := helper.Pagination(offset, limit, total)

	response := response.ConvertToUserGetAllMedicinesResponse(medicine)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"medicines", response, pagination))
}

// User Get Medicine by ID
func GetMedicineUserByIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("medicine_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var medicine schema.Medicine

	if err := configs.DB.First(&medicine, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToUserMedicineResponse(&medicine)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"medicine", response))
}
