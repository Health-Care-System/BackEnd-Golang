package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Create Doctor Transaction
func CreateDoctorTransactionController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	doctorID, _ := strconv.Atoi(c.QueryParam("doctor_id"))

	var doctorTransactionRequest web.CreateDoctorTransactionRequest

	if err := c.Bind(&doctorTransactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input doctor transaction data"))
	}

	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("payment_confirmation")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("payment confirmation file is required"))
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
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid image file format. supported formats: jpg, jpeg, png"))
	}

	paymentConfirmations, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to cloud storage"))
	}

	doctorTransactionRequest.PaymentConfirmation = paymentConfirmations

	if err := helper.ValidateStruct(doctorTransactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var doctor schema.Doctor

	if err := configs.DB.First(&doctor, doctorID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
	}

	doctorTransaction := request.ConvertToCreateDoctorTransactionRequest(doctorTransactionRequest, uint(userID), uint(doctorID), doctor.Fullname, doctor.Specialist, doctor.Price)

	if err := configs.DB.Create(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to create doctor transaction"))
	}

	response := response.ConvertToCreateDoctorTransactionResponse(doctorTransaction, doctor)

	return c.JSON(http.StatusCreated, helper.SuccessResponse("doctor transaction created successful", response))
}

// Get All Doctor Transactions or Get Doctor Transaction by ID or Get Doctor Transaction by Status 
func GetDoctorTransactionsController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	transactionID, _ := strconv.Atoi(c.QueryParam("transaction_id"))
	paymentStatus := c.QueryParam("payment_status")

	if transactionID == 0 && paymentStatus == "" {
		var doctorTransaction []schema.DoctorTransaction

		err := configs.DB.Where("deleted_at IS NULL").Find(&doctorTransaction, "user_id=?", userID).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
		}

		var responses []web.DoctorTransactionsResponse
		for i, doctor_id := range doctorTransaction {
			var doctor schema.Doctor
			err := configs.DB.Find(&doctor, "id=?", doctor_id.DoctorID).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor data"))
			}

			if len(doctorTransaction) == 0 {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("empty doctor transaction data"))
			}

			responses = append(responses, response.ConvertToGetAllDoctorTransactionsResponse(doctorTransaction[i], doctor))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction data successfully retrieved", responses))
	}

	if paymentStatus == "" {
		var doctorTransaction schema.DoctorTransaction
		if err := configs.DB.First(&doctorTransaction, "user_id = ? AND id = ?", userID, transactionID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
		}

		var doctor schema.Doctor
		if err := configs.DB.First(&doctor, "id = ?", doctorTransaction.DoctorID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
		}

		response := response.ConvertToGetDoctorTransactionResponse(doctorTransaction, doctor)

		return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction data successfully retrieved", response))

	}

	if transactionID == 0 {

		var doctorTransactions []schema.DoctorTransaction
		if err := configs.DB.Find(&doctorTransactions, "user_id = ? AND payment_status = ?", userID, paymentStatus).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
		}

		var responses []web.CreateDoctorTransactionResponse
		for i, doctor_id := range doctorTransactions {
			var doctor schema.Doctor
			err := configs.DB.Find(&doctor, "id=?", doctor_id.DoctorID).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
			}

			if len(doctorTransactions) == 0 {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("empty doctor transaction data"))
			}

			responses = append(responses, response.ConvertToGetDoctorTransactionResponse(doctorTransactions[i], doctor))
		}
		return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction data successfully retrieved", responses))

	}
	return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
}

// User Get Doctor Transaction Details by ID
func GetDoctorTransactionDetailsByUserController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	transactionID, _ := strconv.Atoi(c.QueryParam("transaction_id"))

	var doctorTransaction schema.DoctorTransaction

	if err := configs.DB.Preload("Complaint").Preload("Advice").Where("user_id = ? AND id = ?", userID, transactionID).First(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
	}

	response := response.ConvertToDoctorTransactionDetailsResponse(&doctorTransaction)

	return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction details data successfully retrieved", response))
}


// Doctor Get Doctor Transaction Details by ID
func GetDoctorTransactionDetailsByDoctorController(c echo.Context) error {

	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid doctor id"))
	}

	transactionID, _ := strconv.Atoi(c.QueryParam("transaction_id"))

	var doctorTransaction schema.DoctorTransaction

	if err := configs.DB.Preload("Complaint").Preload("Advice").Where(" doctor_id ? AND id = ?", doctorID, transactionID).First(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
	}

	response := response.ConvertToDoctorTransactionDetailsResponse(&doctorTransaction)

	return c.JSON(http.StatusOK, helper.SuccessResponse("doctor transaction details data successfully retrieved", response))
}