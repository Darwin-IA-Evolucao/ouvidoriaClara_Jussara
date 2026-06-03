package controllers

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/usecases"
	"errors"
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
	var msg models.Message
	if err := c.BindJSON(&msg); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	resultado, err := ctrl.usecase.AddMensagem(msg)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, resultado)
}

func (ctrl *MensagemController) GetMessagesByTelefone(c *gin.Context) {
	telefone := c.Param("telefone")
	resposta, err := ctrl.usecase.GetMensagens(telefone)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, resposta)
}

func (ctrl *MensagemController) ClearMessagesByTelefone(c *gin.Context) {
	telefone := c.Param("telefone")
	if err := ctrl.usecase.ClearMensagens(telefone); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Sucesso!")
}
