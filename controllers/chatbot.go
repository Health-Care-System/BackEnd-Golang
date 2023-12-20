package controllers

import (
	"context"
	"fmt"
	"healthcare/models/web"
	"healthcare/utils/helper"
	"healthcare/utils/helper/constanta"
	"healthcare/utils/helper/prompt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

func Chatbot(c echo.Context) error {
	var request web.ChatbotRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	client := openai.NewClient("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")

	response, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("%s %s", prompt.Prompt, request.Request),
			},
		},
	})

	// fmt.Println(string(promptContent), request.Request) // for debugging

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"recommendation"))
	}

	if response.Choices == nil || len(response.Choices) == 0 {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(constanta.ErrActionGet+"recommendation"))
	}

	aiResponse := response.Choices[0].Message.Content
	if strings.Contains(strings.ToLower(aiResponse), "maaf") {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponseWithData(constanta.ErrActionGet+"recommendation", aiResponse))
	} else {
		return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"recommendation", aiResponse))
	}
}

func CustomerService(c echo.Context) error {
	var request web.ChatbotRequest
	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(constanta.ErrInvalidBody))
	}

	if err := helper.ValidateStruct(&request); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	if strings.Contains(strings.ToLower(request.Request), "pembayaran obat") {
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"recommendation", prompt.Content1))

	} else if strings.Contains(strings.ToLower(request.Request), "rating dokter") {
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"recommendation", prompt.Content2))

	} else if strings.Contains(strings.ToLower(request.Request), "riwayat konsultasi") {
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusOK, helper.SuccessResponse(constanta.SuccessActionGet+"recommendation", prompt.Content3))
	} else {
		time.Sleep(10 * time.Second)
		return c.JSON(http.StatusBadRequest, helper.ErrorResponseWithData(constanta.ErrActionGet+"recommendation", "sorry we can't provide answer for your question"))
	}

}
