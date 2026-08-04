package controllers

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/usecases"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampanhaController struct {
	uc usecases.CampanhaUseCase
}

func NewCampanhaController(uc usecases.CampanhaUseCase) *CampanhaController {
	return &CampanhaController{uc: uc}
}

func (ctrl *CampanhaController) CreateCampanha(c *gin.Context) {
	var campanha models.Campanha
	if err := c.ShouldBindJSON(&campanha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.uc.CreateCampanha(&campanha)
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			c.JSON(appError.StatusCode, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, campanha)
}

func (ctrl *CampanhaController) GetCampanhaByPalavraChave(c *gin.Context) {
	palavraChave := c.Param("palavraChave")
	campanha, err := ctrl.uc.GetCampanhaByPalavraChave(palavraChave)
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			c.JSON(appError.StatusCode, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campanha)
}

func (ctrl *CampanhaController) GetAllCampanhas(c *gin.Context) {
	campanhas, err := ctrl.uc.GetAllCampanhas()
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			c.JSON(appError.StatusCode, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campanhas)
}

func (ctrl *CampanhaController) DeleteCampanha(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	err = ctrl.uc.DeleteCampanha(id)
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			c.JSON(appError.StatusCode, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Campanha deletada com sucesso"})
}

func (ctrl *CampanhaController) UpdateCampanha(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	var campanha models.Campanha
	if err := c.ShouldBindJSON(&campanha); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campanha.ID = id
	err = ctrl.uc.UpdateCampanha(&campanha)
	if err != nil {
		var appError *apperror.AppError
		if errors.As(err, &appError) {
			c.JSON(appError.StatusCode, gin.H{"error": appError.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campanha)
}
