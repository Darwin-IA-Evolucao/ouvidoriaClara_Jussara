package controllers

import (
	"back-end/config"
	"back-end/models"
	"back-end/usecases"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LeadController struct {
	useCase usecases.LeadUseCases
}

func NewLeadController(usecase usecases.LeadUseCases) LeadController {
	return LeadController{
		useCase: usecase,
	}
}

func (controller LeadController) DesativarLead(c *gin.Context) {
	telefone := c.Param("telefone")

	res, err := controller.useCase.DesativarLead(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "lead não encontrado"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead desativado com sucesso!"})
}

func (controller LeadController) AtivarLead(c *gin.Context) {
	telefone := c.Param("telefone")

	res, err := controller.useCase.AtivarLead(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "lead não encontrado"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lead reativado com sucesso!"})
}

func (controller LeadController) GetAllLeads(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limite invalido"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset invalido"})
		return
	}
	leads, err := controller.useCase.GetAllLeads(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//estatisticas
	ativos, err := controller.useCase.GetCountContatosAtivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalContatos, err := controller.useCase.GetCountContatos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	planoTotal := config.GetPlanoAtual()

	c.JSON(http.StatusOK, gin.H{
		"leads":    leads,
		"ocupacao": fmt.Sprintf("%d/%d", ativos, planoTotal),
		"limite":   planoTotal,
		"usados":   ativos,
		"total":    totalContatos,
	})
}

func parseOptionalBool(s string) (*bool, error) {
	if s == "" {
		return nil, nil
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (controller LeadController) GetAllContatosUnificados(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limite invalido"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset invalido"})
		return
	}
	filtro := models.ContatosUnificadosFiltro{Search: strings.TrimSpace(c.Query("search"))}
	filtro.Darwin, err = parseOptionalBool(c.Query("darwin"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "darwin inválido"})
		return
	}
	filtro.LeadAtivo, err = parseOptionalBool(c.Query("leadAtivo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "leadAtivo inválido"})
		return
	}
	filtro.Gelo, err = parseOptionalBool(c.Query("gelo"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gelo inválido"})
		return
	}
	filtro.Reclamacao, err = parseOptionalBool(c.Query("reclamacao"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reclamacao inválido"})
		return
	}
	if inicioStr := c.Query("inicio"); inicioStr != "" {
		t, err := time.Parse("2006-01-02", inicioStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "inicio inválido"})
			return
		}
		filtro.Inicio = &t
	}
	if fimStr := c.Query("fim"); fimStr != "" {
		t, err := time.Parse("2006-01-02", fimStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fim inválido"})
			return
		}
		t = t.Add(24 * time.Hour)
		filtro.Fim = &t
	}
	leadsUnificados, err := controller.useCase.GetAllContatosUnificados(limit, offset, filtro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leadsUnificados)
}
