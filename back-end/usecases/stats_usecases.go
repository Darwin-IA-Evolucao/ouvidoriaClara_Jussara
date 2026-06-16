package usecases

import (
	"back-end/models"
	"back-end/repository"
)

type StatsUseCases struct {
	repository repository.StatsRepository
}

func NewStatsUseCases(repo repository.StatsRepository) StatsUseCases {
	return StatsUseCases{repository: repo}
}

func (uc StatsUseCases) GetStats() (*models.Stat, error) {
	pessoasComReclamacao, numConversas, err := uc.repository.GetPessoas()
	if err != nil {
		return nil, err
	}
	numPessoas := pessoasComReclamacao + numConversas

	statsTipo, err := uc.repository.GetCountByTipo()
	if err != nil {
		return nil, err
	}
	statsRegiao, err := uc.repository.GetCountByRegiao()
	if err != nil {
		return nil, err
	}
	statsCategoria, err := uc.repository.GetCountByCategoria()
	if err != nil {
		return nil, err
	}

	stats, err := uc.repository.GetCountByTipoAndStatus()
	if err != nil {
		return nil, err
	}

	percIndicacao := 0.0
	if stats.TotalReclamacoes > 0 {
		percIndicacao = float64(stats.TotalIndicacoes) / float64(stats.TotalReclamacoes) * 100
	}

	percRequerimento := 0.0
	if stats.TotalReclamacoes > 0 {
		percRequerimento = float64(stats.TotalRequerimentos) / float64(stats.TotalReclamacoes) * 100
	}
	
	stat := models.Stat{
		NumPessoas:           int(numPessoas),
		PercIndicacao:        percIndicacao,
		PercRequerimento:     percRequerimento,
		Regioes:              statsRegiao,
		Tipos:                statsTipo,
		Categorias:           statsCategoria,
		StatsByTipoAndStatus: stats,
	}

	return &stat, nil
}
