package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupEnderecoRoutes(router *gin.Engine, enderecoController controllers.EnderecoController) {
	router.GET("/getRegiao", enderecoController.GetRegiao)
	router.POST("/cadastrarEnderecos", enderecoController.CadastrarEnderecos)
}
