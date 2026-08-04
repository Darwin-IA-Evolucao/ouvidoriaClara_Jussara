package usecases

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/repository"
	"fmt"
)

type CampanhaUseCase struct {
	repo *repository.CampanhaRepository
}

func NewCampanhaUseCase(repo *repository.CampanhaRepository) *CampanhaUseCase {
	return &CampanhaUseCase{repo: repo}
}

func (uc *CampanhaUseCase) CreateCampanha(campanha *models.Campanha) error {
	campanha.PalavraChave = campanha.NormalizePalavraChave()
	exists, err := uc.repo.ExistsCampanha(campanha.PalavraChave)
	if err != nil {
		return err
	}
	if exists {
		return apperror.Conflict(fmt.Sprintf("ja existe campanha com a palavra chave %s", campanha.PalavraChave))
	}
	return uc.repo.CreateCampanha(campanha)
}

func (uc *CampanhaUseCase) GetCampanhaByPalavraChave(palavraChave string) (*models.Campanha, error) {
	return uc.repo.GetCampanhaByPalavraChave(palavraChave)
}

func (uc *CampanhaUseCase) GetAllCampanhas() ([]models.Campanha, error) {
	return uc.repo.GetAllCampanhas()
}

func (uc *CampanhaUseCase) GetContatosByCampanha(id int) ([]models.ContatoCampanha, error) {
	campanha, err := uc.repo.GetCampanhaByID(id)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetContatosByCampanha(campanha.PalavraChave)
}

func (uc *CampanhaUseCase) DeleteCampanha(id int) error {
	return uc.repo.DeleteCampanha(id)
}

func (uc *CampanhaUseCase) ExistsCampanha(palavraChave string) (bool, error) {
	return uc.repo.ExistsCampanha(palavraChave)
}

func (uc *CampanhaUseCase) UpdateCampanha(campanha *models.Campanha) error {
	campanha.PalavraChave = campanha.NormalizePalavraChave()
	return uc.repo.UpdateCampanha(campanha)
}
