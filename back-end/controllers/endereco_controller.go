package controllers

import (
	"back-end/models"
	"back-end/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnderecoController struct {
	useCase usecases.EnderecoUseCases
}

func NewEnderecoController(usecase usecases.EnderecoUseCases) EnderecoController {
	return EnderecoController{useCase: usecase}
}

func (ctrl EnderecoController) GetRegiao(c *gin.Context) {
	input := c.Query("rua")
	regiao := ctrl.useCase.GetRegiao(input)
	c.JSON(http.StatusOK, regiao)
}

func (ctrl EnderecoController) CadastrarEnderecos(c *gin.Context) {
	var data []models.Endereco
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if err := ctrl.useCase.CadastrarEnderecos(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Endereço cadastrado com sucesso!"})
}
