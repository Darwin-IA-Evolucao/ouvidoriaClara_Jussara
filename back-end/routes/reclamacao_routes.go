package routes

import (
	"back-end/controllers"

	"github.com/gin-gonic/gin"
)

func SetupReclamacaoRoutes(router *gin.Engine, reclamacaoController controllers.ReclamacaoController) {
	//router.POST("/send", reclamacaoController.CriarReclamacao)
	//router.PUT("/edit/:id", reclamacaoController.EditarReclamacao)
	router.POST("/aprovar/:id", reclamacaoController.AprovarInquerito)
	router.POST("/analise/:id", reclamacaoController.ColocarEmAnalise)
	router.POST("/criado/:id", reclamacaoController.ColocarComoCriado)
	router.POST("/aprovar/requerimento/:id", reclamacaoController.AprovarRequerimento)
	router.POST("/aprovar/causa-animal/:id", reclamacaoController.AprovarCausaAnimal)
	router.POST("/aprovar/outros/:id", reclamacaoController.AprovarOutros)
	//router.POST("/finalizar/:id", reclamacaoController.FinalizarReclamacao)
	// router.POST("/aprovar/oficio/:id", reclamacaoController.AprovarOficio)
	// router.POST("/indicreq/:id", reclamacaoController.AprovarComoAmbos)
	router.POST("/reprovar/:id", reclamacaoController.ReprovarInquerito)
	//router.GET("/reclamacoes", reclamacaoController.GetAllReclamacoes)

	router.POST("/ocorrencia", reclamacaoController.CreateOcorrencia)
	router.GET("/ocorrencias", reclamacaoController.GetAllOcorrencias)
	router.GET("/ocorrencias/relatorio", reclamacaoController.GetRelatorioSolicitacoes)
	router.GET("/ocorrencia/:id", reclamacaoController.GetOcorrenciaById)
	router.PUT("/ocorrencia/:id", reclamacaoController.UpdateOcorrencia)
	router.DELETE("/ocorrencia/:id", reclamacaoController.DeleteOcorrencia)

	router.GET("/categorias", reclamacaoController.GetCategorias)
}
