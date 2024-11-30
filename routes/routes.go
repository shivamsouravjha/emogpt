package routes

import (
	"emogpt/controllers"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	v1 := r.Group("/api")

	{
		v1.POST("/sendMessage", controllers.GptController.SendMessage)
		v1.GET("/keepServerRunning", controllers.HealthController.IsRunning)
	}
}
