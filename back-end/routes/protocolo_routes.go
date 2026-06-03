package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupProtocoloRoutes(router *gin.Engine, protocoloController controllers.ProtocoloController) {
	router.PUT("/protocolo", protocoloController.EnviarProtocolo)
}
