package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLeadRoutes(router *gin.Engine, leadController controllers.LeadController) {
	router.GET("/getLeads", leadController.GetLeads)
	router.GET("/getAllLeads", leadController.GetAllLeads)
	router.DELETE("/desativarLead/:telefone", leadController.DesativarLead)
}
