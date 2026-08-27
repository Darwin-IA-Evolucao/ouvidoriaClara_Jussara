package routes

import (
	"back-end/controllers"
	"back-end/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupContactRoutes(router *gin.Engine, contactController controllers.ContatoController) {
	leadSiteLimiter := middleware.NewRateLimiter(rate.Every(12*time.Second), 3)
	router.POST("/public/lead", middleware.RateLimitPorIP(leadSiteLimiter, 60), contactController.CreateLeadSite)

	router.PUT("/setConversation", contactController.SetConversationId)
	router.GET("/conversation/:telefone", contactController.GetConversationId)
	router.POST("/contatos/import", contactController.ImportContatos)
	router.POST("/contatos", contactController.CreateContato)
	router.GET("/contatos", contactController.GetAllContatos)
	router.GET("/contatos/:telefone", contactController.GetContatoByTelefone)
}
