package controllers

import (
	"back-end/models"
	"back-end/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MensagemController struct {
	usecase *usecases.MensagemUseCase
}

func NewMensagemController(usecase *usecases.MensagemUseCase) *MensagemController {
	return &MensagemController{usecase: usecase}
}

func (ctrl *MensagemController) AddMensagem(c *gin.Context) {
	var addMensagem models.AddMensagem
	if err := c.ShouldBindJSON(&addMensagem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao receber request em AddMensagem: " + err.Error()})
		return
	}
	err := ctrl.usecase.AddMensagem(&addMensagem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar mensagem: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mensagem adicionada com sucesso"})
}

func (ctrl *MensagemController) GetMessagesByTelefone(c *gin.Context) {
	telefone := c.Param("telefone")
	menssagens, err := ctrl.usecase.GetMensagens(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar mensagens: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": menssagens})
}

func (ctrl *MensagemController) ClearMessagesByTelefone(c *gin.Context) {
	telefone := c.Param("telefone")
	err := ctrl.usecase.DeleteMensagens(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao apagar mensagens: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mensagens apagadas com sucesso"})
}
