package controllers

import (
	"errors"
	"fmt"
	"healthcare/configs"
	"healthcare/middlewares"
	"healthcare/models/schema"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/helper/constanta"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllDoctorPagination(offset int, limit int, queryInput []schema.Doctor) ([]schema.Doctor, int64, error) {

	if offset < 0 || limit < 0 {
		return nil, 0, nil
	}

	queryAll := queryInput
	var total int64

	query := configs.DB.Model(&queryAll)

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

// RegisterDoctorController
func RegisterDoctorByAdminController(c echo.Context) error {
	var doctor web.DoctorRegisterRequest

	if err := c.Bind(&doctor); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid input register data"))
	}
	if err := helper.ValidateStruct(doctor); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	doctorRequest := request.ConvertToDoctorRegisterRequest(doctor)

	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("profile_picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("image file is required"))
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

	imageURL, err := helper.UploadFilesToGCS(c, fileHeader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("error upload image to cloud storage"))
	}

	doctorRequest.ProfilePicture = imageURL
	// Periksa apakah email sudah ada
	if existingDoctor := configs.DB.Where("email = ?", doctorRequest.Email).First(&doctorRequest).Error; existingDoctor == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("email already exist"))
	}

	// Save the plain password before hashing
	plainPassword := doctorRequest.Password

	// Hash kata sandi
	doctorRequest.Password = helper.HashPassword(doctor.Password)

	// Simpan data dokter ke database
	if err := configs.DB.Create(&doctorRequest).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Mengirim email pemberitahuan
	includeCredentials := true
	err = helper.SendNotificationEmail(doctorRequest.Email, doctorRequest.Fullname, "register", "doctor", doctorRequest.Email, plainPassword, includeCredentials, 0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to send verification email"))
	}

	response := response.ConvertToDoctorRegisterResponse(doctorRequest)

	return c.JSON(http.StatusCreated, helper.SuccessResponse("registered successful", response))
}

// Login Doctor
func LoginDoctorController(c echo.Context) error {
	var loginRequest web.DoctorLoginRequest

	if err := c.Bind(&loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(loginRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	var doctor schema.Doctor
	if err := configs.DB.Where("email = ? AND deleted_at IS NULL", loginRequest.Email).First(&doctor).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("email not registered"))
	}

	if err := helper.ComparePassword(doctor.Password, loginRequest.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ErrorResponse("incorrect password"))
	}

	// The rest of your code for generating a token and handling the successful login
	token, err := middlewares.GenerateToken(doctor.ID, doctor.Email, doctor.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to generate jwt: "+err.Error()))
	}

	doctorLoginResponse := response.ConvertToDoctorLoginResponse(&doctor)
	doctorLoginResponse.Token = token

	if doctor.Email != "" {
		notificationType := "login"
		userType := "doctor" // Specify the user type as "doctor"
		if err := helper.SendNotificationEmail(doctor.Email, doctor.Fullname, notificationType, userType, "", "", false, 0); err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to send notification email: "+err.Error()))
		}
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("login successful", doctorLoginResponse))
}

func GetAvailableDoctor(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	var doctors []schema.Doctor
	var total int64

	query := configs.DB.Where("status = ?", true)

	query.Model(&doctors).Count(&total)

	query = query.Limit(limit).Offset(offset)

	err = query.Find(&doctors).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctors"))
	}

	if len(doctors) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToGetAllDoctorResponse(doctors)
	pagination := helper.Pagination(offset, limit, total)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"doctors", response, pagination))
}

func GetSpecializeDoctor(c echo.Context) error {
	specialist := c.QueryParam("specialist")

	if specialist == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("specialist"+constanta.ErrQueryParamRequired))
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	var doctors []schema.Doctor
	var total int64

	query := configs.DB.Where("specialist LIKE ? AND status = ?", "%"+specialist+"%", true)

	query.Model(&doctors).Count(&total)

	query = query.Limit(limit).Offset(offset)

	err = query.Find(&doctors).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctors"))
	}

	if len(doctors) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToGetAllDoctorResponse(doctors)
	pagination := helper.Pagination(offset, limit, total)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"doctors", response, pagination))
}

// Get Doctor Profile
func GetDoctorProfileController(c echo.Context) error {
	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	var doctor schema.Doctor
	if err := configs.DB.First(&doctor, doctorID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor profile"))
	}

	response := response.ConvertToDoctorUpdateResponse(&doctor)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"doctor profile", response))
}

// Get All Doctors
func GetAllDoctorByAdminController(c echo.Context) error {

	params := c.QueryParams()
	limit, err := strconv.Atoi(params.Get("limit"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired))
	}

	offset, err := strconv.Atoi(params.Get("offset"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired))
	}

	var doctors []schema.Doctor
	doctor, total, err := GetAllDoctorPagination(offset, limit, doctors)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	pagination := helper.Pagination(offset, limit, total)

	response := response.ConvertToGetAllDoctorByAdminResponse(doctor)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"doctor", response, pagination))
}

// Update Doctor
func UpdateDoctorController(c echo.Context) error {
	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	// Mengambil data dokter yang sudah ada
	var existingDoctor schema.Doctor
	result := configs.DB.First(&existingDoctor, doctorID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"doctor profile"+constanta.ErrNotFound))
	}

	// Parse the request body into the DoctorUpdateRequest struct
	var doctorUpdated web.DoctorsUpdateRequest
	if err := c.Bind(&doctorUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	// apakah email sudah digunakan oleh dokter lain
	var existingDoctorEmail schema.Doctor
	if existingEmail := configs.DB.Where("email = ? AND deleted_at IS NULL", doctorUpdated.Email).First(&existingDoctorEmail).Error; existingEmail == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse(constanta.ErrActionUpdated+"email already exists"))
	}

	// Validate the request payload
	if err := helper.ValidateStruct(doctorUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionUpdated+err.Error()))
	}

	// Hash the password if provided
	if doctorUpdated.Password != "" {
		doctorUpdated.Password = helper.HashPassword(doctorUpdated.Password)
	}

	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("profile_picture")
	if err == nil {
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
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidImageFormat))
		}

		profilePicture, err := helper.UploadFilesToGCS(c, fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+constanta.ErrImageFileRequired))
		}

		doctorUpdated.ProfilePicture = profilePicture
	}

	// Update the doctor details
	if err := configs.DB.Model(&existingDoctor).Updates(doctorUpdated).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+constanta.ErrNotFound))
	}

	configs.DB.Save(&existingDoctor)

	response := response.ConvertToDoctorsUpdateResponse(&existingDoctor)
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"doctor profile", response))
}

// Update Doctor by Admin
func UpdateDoctorByAdminController(c echo.Context) error {
	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid doctor id"))
	}

	// Mengambil data dokter yang sudah ada
	var existingDoctor schema.Doctor
	result := configs.DB.First(&existingDoctor, doctorID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"doctor profile"+constanta.ErrNotFound))
	}

	// Parse the request body into the DoctorUpdateRequest struct
	var doctorUpdated web.DoctorUpdateRequest
	if err := c.Bind(&doctorUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	// apakah email sudah digunakan oleh dokter lain
	var existingDoctorEmail schema.Doctor
	if existingEmail := configs.DB.Where("email = ? AND deleted_at IS NULL", doctorUpdated.Email).First(&existingDoctorEmail).Error; existingEmail == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse(constanta.ErrActionUpdated+"email already exists"))
	}

	// Validate the request payload
	if err := helper.ValidateStruct(doctorUpdated); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionUpdated+err.Error()))
	}

	// Hash the password if provided
	if doctorUpdated.Password != "" {
		doctorUpdated.Password = helper.HashPassword(doctorUpdated.Password)
	}

	err = c.Request().ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	file, fileHeader, err := c.Request().FormFile("profile_picture")
	if err == nil {
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
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidImageFormat))
		}

		profilePicture, err := helper.UploadFilesToGCS(c, fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+constanta.ErrImageFileRequired))
		}

		doctorUpdated.ProfilePicture = profilePicture
	}

	// Update the doctor details
	if err := configs.DB.Model(&existingDoctor).Updates(doctorUpdated).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+constanta.ErrNotFound))
	}

	configs.DB.Save(&existingDoctor)

	response := response.ConvertToDoctorUpdateResponse(&existingDoctor)
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"doctor profile", response))
}

// Delete Doctor
func DeleteDoctorController(c echo.Context) error {
	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	var existingDoctor schema.Doctor
	result := configs.DB.First(&existingDoctor, doctorID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"doctor's account"+constanta.ErrNotFound))
	}

	if err := configs.DB.Delete(&existingDoctor).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionDeleted+"doctor's account"+constanta.ErrNotFound))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionDeleted+"doctor account", nil))
}

// DeleteDoctorByAdminController deletes a doctor by admin
func DeleteDoctorByAdminController(c echo.Context) error {
	// Parse doctor ID from the request parameters
	doctor_id, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid id"))
	}

	// Retrieve the existing doctor from the database
	var existingDoctor schema.Doctor
	result := configs.DB.First(&existingDoctor, doctor_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve data"))
	}

	// Delete the doctor from the database
	result = configs.DB.Delete(&existingDoctor, doctor_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to delete data"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse("data deleted successfuly  ", nil))
}

// Get Doctor by ID
func GetDoctorByIDController(c echo.Context) error {
	doctorID, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	var doctor schema.Doctor
	result := configs.DB.First(&doctor, doctorID)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	response := response.ConvertToGetIDDoctorResponse(&doctor)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"doctor details", response))
}

// Get Doctor by ID
func GetDoctorIDbyAdminController(c echo.Context) error {
	doctor_id, err := strconv.Atoi(c.Param("doctor_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("failed to retrieve doctor id"))
	}

	var doctor schema.Doctor
	result := configs.DB.First(&doctor, doctor_id)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to fetch doctor data"))
	}

	response := response.ConvertToGetDoctorbyAdminResponse(&doctor)

	return c.JSON(http.StatusOK, helper.SuccessResponse("data successfully retrieved", response))
}

// Manage User
func GetManageUserController(c echo.Context) error {
	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse((constanta.ErrActionGet + "doctor id")))
	}

	// Parse limit and offset from query parameters
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("limit"+constanta.ErrQueryParamRequired+": "+err.Error()))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("offset"+constanta.ErrQueryParamRequired+": "+err.Error()))
	}

	transactionID, _ := strconv.Atoi(c.QueryParam("transaction_id"))
	patientStatus := c.QueryParam("patient_status")
	fullname := c.QueryParam("fullname")
	keyword := c.QueryParam("keyword")

	var manageUser []schema.DoctorTransaction
	var total int64

	var query *gorm.DB

	if keyword != "" {
		// Search by keyword (transactionID, patient_status, or fullname)
		query = configs.DB.
			Where("doctor_transactions.doctor_id = ? AND doctor_transactions.payment_status = 'success' AND (doctor_transactions.id LIKE ? OR doctor_transactions.patient_status LIKE ? OR users.fullname LIKE ?)", doctorID, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
			Joins("JOIN users ON doctor_transactions.user_id = users.id").
			Order("doctor_transactions.created_at desc")
	} else if transactionID != 0 {
		// Get transaction by ID
		query = configs.DB.Where("doctor_transactions.doctor_id = ? AND doctor_transactions.id = ? AND doctor_transactions.payment_status = 'success'", doctorID, transactionID)
	} else if patientStatus != "" {
		// Get transactions by patient status
		query = configs.DB.Where("doctor_transactions.doctor_id = ? AND doctor_transactions.patient_status LIKE ? AND doctor_transactions.payment_status = 'success'", doctorID, "%"+patientStatus+"%").Order("doctor_transactions.created_at desc")
	} else if fullname != "" {
		// Get transactions by fullname
		query = configs.DB.
			Joins("JOIN users ON doctor_transactions.user_id = users.id").
			Where("doctor_transactions.doctor_id = ? AND users.fullname LIKE ? AND doctor_transactions.payment_status = 'success'", doctorID, "%"+fullname+"%").Order("doctor_transactions.created_at desc")
	} else {
		// Get all transactions
		query = configs.DB.Where("doctor_transactions.deleted_at IS NULL AND doctor_transactions.payment_status = 'success'").Where("doctor_transactions.doctor_id=?", doctorID).Order("doctor_transactions.created_at desc")
	}

	// Count total number of records
	if err := query.Model(&schema.DoctorTransaction{}).Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"count doctor transaction: "+err.Error()))
	}

	// Apply limit and offset to the query
	query = query.Limit(limit).Offset(offset)

	// Execute the query
	if err := query.Find(&manageUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor transaction: "+err.Error()))
	}

	if len(manageUser) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(fmt.Sprintf("no doctor transaction data found for doctorid: %d, transactionID: %d, patientStatus: %s", doctorID, transactionID, patientStatus)))
	}

	var responses []web.ManageUserResponse
	for _, doctorTransaction := range manageUser {
		var user schema.User
		userID := doctorTransaction.UserID
		if err := configs.DB.First(&user, "id=?", userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"user: "+err.Error()))
			}
		}

		response := response.ConvertToManageUserResponse(doctorTransaction, user)
		responses = append(responses, response)
	}

	pagination := helper.Pagination(offset, limit, total)

	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"doctor transaction", responses, pagination))
}

// Update manage user
func UpdateManageUserController(c echo.Context) error {

	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse((constanta.ErrActionGet + "doctor id")))
	}

	// Mendapatkan data permintaan
	var requestBody web.UpdateManageUserRequest
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}
	if err := helper.ValidateStruct(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Memeriksa Update
	if requestBody.HealthDetails == "" && requestBody.PatientStatus == "" {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("health details or patient status is required"))
	}

	// Mengambil ID transaksi dari parameter
	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidIDParam))
	}

	// Mengambil data transaksi dokter berdasarkan ID dokter dan ID transaksi
	var doctorTransaction schema.DoctorTransaction
	err = configs.DB.First(&doctorTransaction, "doctor_id = ? AND id = ?", doctorID, transactionID).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound))
	}

	// Memeriksa apakah status pembayaran adalah "success"
	if doctorTransaction.PaymentStatus != "success" {
		return c.JSON(http.StatusForbidden, helper.ErrorResponse("payment status is not 'success'"))
	}

	// Memperbarui Update
	if requestBody.HealthDetails != "" {
		doctorTransaction.HealthDetails = requestBody.HealthDetails
	}

	if requestBody.PatientStatus != "" {
		doctorTransaction.PatientStatus = requestBody.PatientStatus
	}

	// Menyimpan transaksi dokter yang diperbarui ke database
	if err := configs.DB.Save(&doctorTransaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"health details and patient status"))
	}

	// Mendapatkan data pengguna
	var user schema.User
	err = configs.DB.First(&user, "id=?", doctorTransaction.UserID).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"user data"))
	}

	response := response.ConvertToManageUserResponse(doctorTransaction, user)
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"health details and patient status", response))
}

func GetAllDoctorConsultationController(c echo.Context) error {
	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid limit parameter"))
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid offset parameter"))
	}

	var consultations []schema.DoctorTransaction
	var total int64

	// Bangun query database untuk mengambil konsultasi dokter
	query := configs.DB.
		Joins("INNER JOIN roomchats ON doctor_transactions.id = roomchats.transaction_id").
		Joins("LEFT JOIN messages ON roomchats.id = messages.roomchat_id").
		Where("doctor_transactions.payment_status = ?", "success").
		Where("roomchats.transaction_id IS NOT NULL").
		Where("doctor_transactions.doctor_id = ?", doctorID).
		Where("messages.ID IS NULL").
		Order("doctor_transactions.created_at DESC")

	query.Model(&consultations).Count(&total)

	query = query.Limit(limit).Offset(offset)

	// Ambil konsultasi berdasarkan query akhir
	if err := query.Find(&consultations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"consultations"))
	}

	var consultationResponses []web.DoctorConsultationResponse
	for _, consultation := range consultations {
		var user schema.User
		if err := configs.DB.First(&user, consultation.UserID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"user"))
		}

		// Use Preload to fetch the associated Roomchat
		if err := configs.DB.Preload("Roomchat").First(&consultation, consultation.ID).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"room"))
		}

		response := response.ConvertToConsultationResponse(consultation, user, consultation.Roomchat)
		consultationResponses = append(consultationResponses, response)
	}

	if len(consultationResponses) == 0 {
		return c.JSON(http.StatusNotFound, helper.ErrorResponse(constanta.ErrNotFound+"consultations"))
	}

	pagination := helper.Pagination(offset, limit, total)
	return c.JSON(http.StatusOK, helper.PaginationResponse(constanta.SuccessActionGet+"consultations", consultationResponses, pagination))
}

// Change Doctor Status
func ChangeDoctorStatusController(c echo.Context) error {

	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	// Parse the request body
	var statusRequest web.ChangeDoctorStatusRequest
	if err := c.Bind(&statusRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	// Validate input status
	if err := helper.ValidateStruct(statusRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Retrieve existing doctor data
	var existingDoctor schema.Doctor
	result := configs.DB.First(&existingDoctor, doctorID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+"doctor profile"+constanta.ErrNotFound))
	}

	// Update doctor status
	existingDoctor.Status = statusRequest.Status
	if err := configs.DB.Save(&existingDoctor).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionUpdated+constanta.ErrNotFound))
	}

	response := response.ConvertToDoctorUpdateResponse(&existingDoctor)
	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"doctor status", response))
}

// reset password dan mengirimkan OTP ke email
func ResetPasswordDoctor(c echo.Context) error {
	var resetRequest web.ResetRequest
	if err := c.Bind(&resetRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request"))
	}

	if err := helper.ValidateStruct(resetRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Verify OTP
	if err := helper.VerifyOTPByEmail(resetRequest.Email, resetRequest.OTP, "doctor"); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionGet+"OTP verification failed"))
	}

	hashedPassword := helper.HashPassword(resetRequest.Password)

	// Update password
	if err := helper.UpdatePasswordInDatabase(configs.DB, "doctors", resetRequest.Email, hashedPassword, resetRequest.OTP); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"update password"))
	}

	// Delete OTP from the database
	if err := helper.DeleteOTPFromDatabase(configs.DB, "doctors", resetRequest.Email); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"delete OTP"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionUpdated+"doctor's password", nil))
}

func GetOTPForPasswordDoctor(c echo.Context) error {
	var OTPRequest web.PasswordResetRequest
	if err := c.Bind(&OTPRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(OTPRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	if err := helper.SendOTPViaEmail(OTPRequest.Email, "doctor", "reset"); err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"send OTP"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionCreated+"OTP", nil))
}

func VerifyOTPDoctor(c echo.Context) error {
	var verificationRequest web.OTPVerificationRequest
	if err := c.Bind(&verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid request"))
	}

	if err := helper.ValidateStruct(verificationRequest); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	// Verify OTP and handle errors
	if err := helper.VerifyOTPByEmail(verificationRequest.Email, verificationRequest.OTP, "doctor"); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrActionGet+"OTP not found"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"OTP verification", nil))
}

func GetDoctorStatusController(c echo.Context) error {

	doctorID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor id"))
	}

	// Mengambil data dokter dari database
	var doctor schema.Doctor
	result := configs.DB.Model(&schema.Doctor{}).Where("id = ?", doctorID).First(&doctor)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("doctor's data"+constanta.ErrNotFound))
		}
		// Menghandle error lain sesuai kebutuhan
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"doctor's data"))
	}

	response := response.ConvertToDoctorStatusResponse(&doctor)

	return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"doctor status", response))
}
