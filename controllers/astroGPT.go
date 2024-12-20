package controllers

import (
	"emogpt/services"
	"emogpt/utils/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AstroControllerI interface {
	SendMessage(ctx *gin.Context)
}

type astroController struct{}

var AstroController AstroControllerI = &astroController{}

func (s *astroController) SendMessage(ctx *gin.Context) {
	var requestBody services.AstroRequestBody

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	promptBuilder, _ := services.NewAstroPromptBuilder(requestBody)

	prompt, _ := promptBuilder.BuildAstroPrompt("meta")
	av, _ := helpers.GenerateChat(ctx, *prompt)
	ctx.JSON(200, gin.H{"message": av})
}
