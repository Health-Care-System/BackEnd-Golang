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
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllDoctorTransactionPagination(userID int, offset int, limit int, paymentStatus string, queryInput []schema.DoctorTransaction) ([]schema.DoctorTransaction, int64, error) {

	if offset < 0 || limit < 0 {
		return nil, 0, nil
	}

	queryAll := queryInput
	var total int64

	query := configs.DB.Model(&queryAll)

	if paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
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

// Create Doctor Transaction
func CreateDoctorTransactionController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid doctor id"))
	}

	var doctorTransactionRequest web.CreateDoctorTransactionRequest

	if err := c.Bind(&doctorTransactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input doctor transaction data"))
	}

	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("payment_confirmation")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("payment confirmation file is required"))
	}
	defer file.Close()

	if fileHeader.Size > 10*1024*1024 { // 10 MB limit
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("image file size exceeds the limit (10 MB)"))
	}

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
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid image file format. supported formats: jpg, jpeg, png"))
	}

	paymentConfirmations, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to cloud storage"))
	}

	doctorTransactionRequest.PaymentConfirmation = paymentConfirmations

	if err := helper.ValidateStruct(doctorTransactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	paymentMethod := doctorTransactionRequest.PaymentMethod

	if !helper.PaymentMethodIsValid(paymentMethod) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input payment method data ('manual transfer bca', 'manual transfer bri', 'manual transfer bni')"))
	}

	var doctor schema.Doctor

	if err := configs.DB.First(&doctor, "id = ? AND status = true", doctorID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
	}

	doctorTransaction := request.ConvertToCreateDoctorTransactionRequest(doctorTransactionRequest, uint(userID), uint(doctorID), doctor.Fullname, doctor.Specialist, doctor.Price)

	if err := configs.DB.Create(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to create doctor transaction"))
	}

	response := response.ConvertToCreateDoctorTransactionResponse(doctorTransaction, doctor)

	return c.JSON(http.StatusCreated, helper.SuccessResponse("doctor transaction created successful", response))
}

// Get Doctor Transaction by ID
func GetDoctorTransactionController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		log.Println(ok)
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid transaction id"))
	}

	var doctorTransaction schema.DoctorTransaction

	if err := configs.DB.First(&doctorTransaction, "user_id = ? AND id = ?", userID, transactionID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
	}

	var doctor schema.Doctor
	if err := configs.DB.First(&doctor, "id = ?", doctorTransaction.DoctorID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
	}

	response := response.ConvertToGetDoctorTransactionResponse(doctorTransaction, doctor)

	return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction data successfully retrieved", response))
}

// Get All Doctor Transactions or Get Doctor Transaction by Status
func GetAllDoctorTransactionsController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	var total int64

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	paymentStatus := c.QueryParam("payment_status")

	if !helper.PaymentStatusIsValid(paymentStatus) {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input payment status data ('pending', 'success', 'cancelled')"))
	}

	// Get All Doctor Transactions
	if paymentStatus == "" {

		var doctorTransaction []schema.DoctorTransaction

		doctorTransactions, total, err := GetAllDoctorTransactionPagination(userID, offset, limit, paymentStatus, doctorTransaction)
		if err != nil {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("doctor transaction data not found"))
		}

		err = configs.DB.Find(&doctorTransaction, "user_id=?", userID).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
		}

		var responses []web.DoctorTransactionsResponse
		for _, doctorTransaction := range doctorTransactions {

			var doctor schema.Doctor
			err := configs.DB.Find(&doctor, "id=?", doctorTransaction.DoctorID).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
			}

			responses = append(responses, response.ConvertToGetAllDoctorTransactionsResponse(doctorTransaction, doctor))
		}

		if len(doctorTransaction) == 0 {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("empty doctor transaction data"))
		}

		pagination := helper.Pagination(offset, limit, total)

		return c.JSON(http.StatusOK, helper.PaginationResponse("doctor transaction data successfully retrieved", responses, pagination))
	}

	// Get Doctor Transactions by Status

	var doctorTransaction []schema.DoctorTransaction

	doctorTransactions, total, err := GetAllDoctorTransactionPagination(userID, offset, limit, paymentStatus, doctorTransaction)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
	}

	var responses []web.CreateDoctorTransactionResponse
	for _, doctorTransaction := range doctorTransactions {

		var doctor schema.Doctor
		err := configs.DB.Find(&doctor, "id=?", doctorTransaction.DoctorID).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
		}

		responses = append(responses, response.ConvertToGetDoctorTransactionResponse(doctorTransaction, doctor))
	}

	if len(doctorTransactions) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse("empty doctor transaction data"))
	}

	pagination := helper.Pagination(offset, limit, total)

	return c.JSON(http.StatusOK, helper.PaginationResponse("doctor transaction data successfully retrieved", responses, pagination))
}
