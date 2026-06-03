package controllers

import (
	"back-end/config"
	"back-end/usecases"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

type LeadController struct {
	useCase usecases.LeadUseCases
}

func NewLeadController (usecase usecases.LeadUseCases) LeadController{
	return LeadController{
		useCase: usecase,
	}
}

func (controller LeadController)  DesativarLead(c *gin.Context){
	telefone := c.Param("telefone")

	res,err := controller.useCase.DesativarLead(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected,err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0{
		c.JSON(http.StatusNotFound, gin.H{"error":"lead não encontrado"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead desativado com sucesso!"})
}

func (controller LeadController)  AtivarLead(c *gin.Context){
	telefone := c.Param("telefone")

	res,err := controller.useCase.AtivarLead(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected,err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0{
		c.JSON(http.StatusNotFound, gin.H{"error":"lead não encontrado"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead reativado com sucesso!"})
}

func (controller LeadController) GetAllLeads(c *gin.Context){
	leads, err := controller.useCase.GetAllLeads()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	//estatisticas
	ativos, err := controller.useCase.GetCountContatosAtivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	totalContatos, err := controller.useCase.GetCountContatos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}

	planoTotal := config.GetPlanoAtual()

	c.JSON(http.StatusOK, gin.H{
		"leads":leads,
		"ocupacao": fmt.Sprintf("%d/%d",ativos, planoTotal),
		"limite": planoTotal,
		"usados": ativos,
		"total": totalContatos,
	})
}