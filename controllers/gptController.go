package controllers

import (
	"emogpt/services"
	"emogpt/utils/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GptControllerI interface {
	SendMessage(ctx *gin.Context)
}

type gptController struct{}

var GptController GptControllerI = &gptController{}

func (s *gptController) SendMessage(ctx *gin.Context) {
	var requestBody struct {
		Message string `json:"message"`
		Mood    string `json:"mood"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}
	message := requestBody.Message

	promptBuilder, _ := services.NewPromptBuilder(message, requestBody.Mood)

	prompt, _ := promptBuilder.BuildPrompt("meta")
	av, _ := helpers.GenerateChat(ctx, *prompt)
	ctx.JSON(200, gin.H{"message": av.Response})
}
