package controllers

import (
	"back-end/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsController struct {
	useCase usecases.StatsUseCases
}

func NewStatsController(usecase usecases.StatsUseCases) StatsController {
	return StatsController{useCase: usecase}
}

func (ctrl StatsController) GetStats(c *gin.Context) {
	stats, err := ctrl.useCase.GetStats()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, stats)
}
