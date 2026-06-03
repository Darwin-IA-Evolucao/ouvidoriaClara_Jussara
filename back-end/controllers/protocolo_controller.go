package controllers

import (
	"back-end/models"
	"back-end/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProtocoloController struct {
	useCase usecases.ProtocoloUseCases
}

func NewProtocoloController(usecase usecases.ProtocoloUseCases) ProtocoloController {
	return ProtocoloController{useCase: usecase}
}

func (ctrl ProtocoloController) EnviarProtocolo(c *gin.Context) {
	var protocolo models.Protocolo
	if err := c.BindJSON(&protocolo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if err := ctrl.useCase.EnviarProtocolo(protocolo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Protocolo enviado e cliente notificado com sucesso."})
}
