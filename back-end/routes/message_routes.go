package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetMessageRoutes(router *gin.Engine, messageController *controllers.MensagemController) {
	router.POST("/addMessage", messageController.AddMensagem)
	router.GET("/getMessages/:telefone", messageController.GetMessagesByTelefone)
	router.DELETE("/clearMessages/:telefone", messageController.ClearMessagesByTelefone)
}
