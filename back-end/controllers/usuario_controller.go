package controllers

import (
	"back-end/models"
	"back-end/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UsuarioController struct {
	useCase usecases.UsuarioUseCases
}

func NewUsuarioController(usecase usecases.UsuarioUseCases) UsuarioController {
	return UsuarioController{useCase: usecase}
}

func (controller UsuarioController) Login(c *gin.Context) {
	var usuario models.Usuario
	if err := c.BindJSON(&usuario); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	autenticado, err := controller.useCase.Login(usuario)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !autenticado {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Credenciais inválidas."})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Login sucesso.", "expiration_time": time.Now().Add(24 * time.Hour).Unix()})
}

func (controller UsuarioController) Greetings(c *gin.Context) {
	c.String(http.StatusOK, "Jesus é o Senhor!")
}
