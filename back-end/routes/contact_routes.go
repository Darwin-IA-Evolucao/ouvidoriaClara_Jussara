package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupContactRoutes(router *gin.Engine, contactController controllers.ContatoController) {
	router.PUT("/setConversation", contactController.SetConversationId)
	router.GET("/conversation/:telefone", contactController.GetConversationId)
	router.POST("/contatos/import", contactController.ImportContatos)
	router.POST("/contatos", contactController.CreateContato)
	router.GET("/contatos", contactController.GetAllContatos)
	router.GET("/contatos/:telefone", contactController.GetContatoByTelefone)
}
