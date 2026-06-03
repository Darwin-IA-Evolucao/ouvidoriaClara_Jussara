package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupContactRoutes(router *gin.Engine, contactController controllers.ContatoController) {
	router.PUT("/setConversation", contactController.SetConversationId)
	router.GET("/conversation/:telefone", contactController.GetConversationId)
	router.GET("/contatos", contactController.GetAllContatos)
	router.GET("/contatos/:telefone", contactController.GetContatoByTelefone)
}
