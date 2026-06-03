package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupLeadRoutes(router *gin.Engine, leadController controllers.LeadController) {
	router.DELETE("/leads/:telefone", leadController.DesativarLead)  
	router.PUT("/leads/ativar/:telefone", leadController.AtivarLead) 
	router.GET("/leads", leadController.GetAllLeads)                
}
