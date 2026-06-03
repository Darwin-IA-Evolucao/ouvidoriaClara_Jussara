package controllers

import (
	"back-end/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeadController struct {
	useCase usecases.LeadUseCases
}

func NewLeadController(usecase usecases.LeadUseCases) LeadController {
	return LeadController{useCase: usecase}
}

func (controller LeadController) GetLeads(c *gin.Context) {
	leads, err := controller.useCase.GetLeads()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, leads)
}

func (controller LeadController) GetAllLeads(c *gin.Context) {
	leads, ocupacao, err := controller.useCase.GetAllLeads()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"leads": leads, "ocupacao": ocupacao})
}

func (controller LeadController) DesativarLead(c *gin.Context) {
	telefone := c.Param("telefone")
	if err := controller.useCase.DesativarLead(telefone); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Lead desativado com sucesso.")
}
