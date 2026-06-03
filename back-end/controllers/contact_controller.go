package controllers

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/usecases"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContatoController struct {
	usecase *usecases.ContatoUseCase
}

func NewContatoController(usecase *usecases.ContatoUseCase) *ContatoController {
	return &ContatoController{usecase: usecase}
}

func (ctrl *ContatoController) InitializeContact(c *gin.Context) {
	var contact models.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if contact.Telefone == "status" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Telefone obrigatório."})
		return
	}

	resultado, err := ctrl.usecase.InitializeContact(contact)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resultado == "novo" {
		c.String(http.StatusCreated, "novo")
		return
	}
	c.String(http.StatusOK, resultado)
}

func (ctrl *ContatoController) SetConversationId(c *gin.Context) {
	var mold models.ContactSetMold
	if err := c.BindJSON(&mold); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if err := ctrl.usecase.SetConversationId(mold.Telefone, mold.ConversationID); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Sucesso!")
}
