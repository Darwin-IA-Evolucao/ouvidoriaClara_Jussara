package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetUsuarioRoutes(router *gin.Engine, usuarioController controllers.UsuarioController) {
	router.POST("/login", usuarioController.Login)
	router.GET("/", usuarioController.Greetings)
}
