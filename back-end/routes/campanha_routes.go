package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCampanhaRoutes(router *gin.Engine, ctrl *controllers.CampanhaController) {
	router.POST("/campanhas", ctrl.CreateCampanha)
	router.GET("/campanhas/palavra-chave/:palavraChave", ctrl.GetCampanhaByPalavraChave)
	router.GET("/campanhas/:id/contatos", ctrl.GetContatosByCampanha)
	router.GET("/campanhas", ctrl.GetAllCampanhas)
	router.DELETE("/campanhas/:id", ctrl.DeleteCampanha)
	router.PUT("/campanhas/:id", ctrl.UpdateCampanha)
}
