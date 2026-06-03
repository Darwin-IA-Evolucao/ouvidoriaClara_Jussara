package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupStatsRoutes(router *gin.Engine, statsController controllers.StatsController) {
	router.GET("/stats", statsController.GetStats)
}
