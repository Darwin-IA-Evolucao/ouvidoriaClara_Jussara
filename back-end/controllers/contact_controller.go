package controllers

import (
	"back-end/models"
	"back-end/usecases"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContatoController struct {
	usecase *usecases.ContatoUseCase
}

func NewContatoController(usecase *usecases.ContatoUseCase) *ContatoController {
	return &ContatoController{usecase: usecase}
}

func (ctrl *ContatoController) GetAllContatos(c *gin.Context) {
	contatos, err := ctrl.usecase.GetAllContatos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar contatos: " + err.Error()})
		fmt.Println("erro ao buscar contato: ", err.Error())
		return
	}
	c.JSON(http.StatusOK, contatos)
}

func (ctrl *ContatoController) SetConversationId(c *gin.Context) {
	var request models.SetConversationIdRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao receber request em SetConversationId: " + err.Error()})
		return
	}

	err := ctrl.usecase.SetConversationId(request.Telefone, request.ConversationId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao setar conversationId: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ConversationId setado com sucesso"})
}

func (ctrl *ContatoController) GetConversationId(c *gin.Context) {
	telefone := c.Param("telefone")
	conversation, err := ctrl.usecase.GetConversationIdByTelefone(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar conversationId: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, conversation)
}

func (ctrl *ContatoController) GetContatoByTelefone(c *gin.Context) {
	telefone := c.Param("telefone")
	contato, err := ctrl.usecase.GetContatoByTelefone(telefone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar contato: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, contato)
}
