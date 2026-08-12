package controllers

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/usecases"
	"errors"
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

func (ctrl *ContatoController) CreateContato(c *gin.Context) {
	var request models.CreateContatoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.usecase.CreateContato(request.Telefone, request.Nome, request.Campanha)
	if err != nil {
		respondContatoErr(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "contato criado"})
}

func (ctrl *ContatoController) ImportContatos(c *gin.Context) {
	campanha := c.PostForm("campanha")
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "arquivo obrigatorio"})
		return
	}
	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "erro ao abrir arquivo"})
		return
	}
	defer file.Close()
	result, err := ctrl.usecase.ImportContatos(header.Filename, file, campanha)
	if err != nil {
		respondContatoErr(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}


func respondContatoErr(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
