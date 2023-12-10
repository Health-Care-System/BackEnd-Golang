package controllers

import (
	"healthcare/configs"
	"healthcare/models/schema"
	"healthcare/utils/helper"
	"healthcare/utils/request"
	"healthcare/utils/response"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// User Create Roomchat and Send Notification to Doctor
func CreateRoomchatsController(c echo.Context) error {

	userID, ok := c.Get("userID").(int)
	if !ok {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("invalid user id"))
	}

	transactionID, err := strconv.Atoi(c.Param("transaction_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("invalid transaction id"))
	}

	var existingRoomchat schema.Roomchat
	if err := configs.DB.First(&existingRoomchat, "transaction_id = ?", transactionID).Error; err == nil {
		return c.JSON(http.StatusConflict, helper.ErrorResponse("roomchat for this transaction id already exists"))
	}

	var doctorTransaction schema.DoctorTransaction
	if err := configs.DB.First(&doctorTransaction, "user_id = ? AND id = ? AND payment_status = 'success'", userID, transactionID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to retrieve doctor transaction data"))
	}

	roomchat := request.CreateRoomchatRequest(uint(transactionID))

	if err := configs.DB.Create(&roomchat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to create roomchat"))
	}

	response := response.ConvertToCreateRoomchatResponse(&roomchat)

	return c.JSON(http.StatusCreated, helper.SuccessResponse("roomchat created successful", response))

}

// User or Doctor Join Roomchat
func JoinRoomchatController(c echo.Context) error {
	return nil
}

// User or Doctor Get Roomchat
func GetRoomchatController(c echo.Context) error {
	return nil
}

// Doctor Get All Roomchats
func GetAllRoomchatsController(c echo.Context) error {
	return nil
}

// Doctor Get Client
func GetClientController(c echo.Context) error {
	return nil
}
