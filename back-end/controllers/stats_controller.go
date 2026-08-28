package controllers

import (
	"back-end/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StatsController struct {
	useCase usecases.StatsUseCases
}

func NewStatsController(usecase usecases.StatsUseCases) StatsController {
	return StatsController{useCase: usecase}
}

func (ctrl StatsController) GetStats(c *gin.Context) {
	var inicio, fim *time.Time
	if inicioStr := c.Query("inicio"); inicioStr != "" {
		t, err := time.Parse("2006-01-02", inicioStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "inicio inválido"})
			return
		}
		inicio = &t
	}
	if fimStr := c.Query("fim"); fimStr != "" {
		t, err := time.Parse("2006-01-02", fimStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fim inválido"})
			return
		}
		t = t.Add(24 * time.Hour)
		fim = &t
	}
	stats, err := ctrl.useCase.GetStats(inicio, fim)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, stats)
}
