package controllers

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/usecases"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReclamacaoController struct {
	useCase usecases.ReclamacaoUseCases
}

func NewReclamacaoController(usecase usecases.ReclamacaoUseCases) ReclamacaoController {
	return ReclamacaoController{useCase: usecase}
}

func (ctrl ReclamacaoController) GetCategorias(c *gin.Context) {
	c.JSON(http.StatusOK, ctrl.useCase.GetCategorias())
}

func (ctrl ReclamacaoController) CriarReclamacao(c *gin.Context) {
	var data models.RequestData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos são obrigatórios"})
		return
	}
	if err := ctrl.useCase.CreateReclamacao(data); err != nil {
		if err.Error() == "Categoria inválida" {
			c.JSON(http.StatusOK, gin.H{"error": "Categoria inválida"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Solicitação enviada com sucesso!"})
}

func (ctrl ReclamacaoController) EditarReclamacao(c *gin.Context) {
	id := c.Param("id")
	var data models.Atualizacao
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if err := ctrl.useCase.EditarReclamacao(id, data.NovaReclamacao); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reclamação editada com sucesso!"})
}

func (ctrl ReclamacaoController) GetAllReclamacoes(c *gin.Context) {
	reclamacoes, err := ctrl.useCase.GetAllReclamacoes()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, reclamacoes)
}

func (ctrl ReclamacaoController) AprovarInquerito(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.AprovarInquerito(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito aprovado e enviado com sucesso!"})
}

func (ctrl ReclamacaoController) ColocarEmAnalise(c *gin.Context) {
	// colocar uma data para retomar a reclamação
	var req struct {
		Data      string `json:"data"`
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.ColocarEmAnalise(c.Param("id"), req.Data, req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito aprovado e enviado com sucesso!"})
}
func (ctrl ReclamacaoController) ColocarComoCriado(c *gin.Context) {
	var req struct {
		IDUsuario int `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.ColocarComoCriado(c.Param("id"), req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito aprovado e enviado com sucesso!"})
}

func (ctrl ReclamacaoController) AprovarRequerimento(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.AprovarRequerimento(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Requerimento aprovado e enviado com sucesso!"})
}

// func (ctrl ReclamacaoController) AprovarOficio(c *gin.Context) {
// 	if err := ctrl.useCase.AprovarOficio(c.Param("id")); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Ofício aprovado e enviado com sucesso!"})
// }

func (ctrl ReclamacaoController) AprovarComoAmbos(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.AprovarComoAmbos(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inqueritos aprovados e enviados com sucesso!"})
}

func (ctrl ReclamacaoController) ReprovarInquerito(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.ReprovarInquerito(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito reprovado com sucesso!"})
}

func (ctrl ReclamacaoController) AprovarCausaAnimal(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.AprovarCausaAnimal(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Causa animal aprovada com sucesso!"})
}

func (ctrl ReclamacaoController) FinalizarReclamacao(c *gin.Context) {
	var req struct {
		Mensagem  string `json:"mensagem"`
		IDUsuario int    `json:"idUsuario"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ctrl.useCase.FinalizarReclamacao(c.Param("id"), req.Mensagem, req.IDUsuario); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reclamacao finalizada com sucesso!"})
}

func (ctrl ReclamacaoController) CreateOcorrencia(c *gin.Context) {
	var req models.OcorrenciaRequest
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("erro ao bindar ocorrencia: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	id, err := ctrl.useCase.CreateOcorrencia(req)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			fmt.Println("erro ao criar ocorrencia: ", err)
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		fmt.Println("erro ao criar ocorrencia: ", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Ocorrência registrada com sucesso", "id": id})
}

func (ctrl ReclamacaoController) GetAllOcorrencias(c *gin.Context) {
	telefone := c.Query("telefone")
	list, err := ctrl.useCase.GetAllOcorrencias(telefone)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if list == nil {
		list = []models.Ocorrencia{}
	}
	c.IndentedJSON(http.StatusOK, list)
}

func (ctrl ReclamacaoController) GetRelatorioSolicitacoes(c *gin.Context) {
	var inicio, fim *time.Time
	if inicioStr := c.Query("inicio"); inicioStr != "" {
		t, err := time.Parse("2006-01-02", inicioStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "inicio inválido"})
			return
		}
		inicio = &t
	}
	if fimStr := c.Query("fim"); fimStr != "" {
		t, err := time.Parse("2006-01-02", fimStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fim inválido"})
			return
		}
		t = t.Add(24 * time.Hour)
		fim = &t
	}
	relatorio, err := ctrl.useCase.GetRelatorioSolicitacoes(inicio, fim)
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, relatorio)
}

func (ctrl ReclamacaoController) GetOcorrenciaById(c *gin.Context) {
	o, err := ctrl.useCase.GetOcorrenciaById(c.Param("id"))
	if err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, o)
}

func (ctrl ReclamacaoController) UpdateOcorrencia(c *gin.Context) {
	var req models.OcorrenciaUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	if err := ctrl.useCase.UpdateOcorrencia(c.Param("id"), req); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Ocorrência atualizada com sucesso"})
}

func (ctrl ReclamacaoController) DeleteOcorrencia(c *gin.Context) {
	if err := ctrl.useCase.DeleteOcorrencia(c.Param("id")); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			c.IndentedJSON(appErr.StatusCode, gin.H{"error": appErr.Message})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Ocorrência removida com sucesso"})
}
