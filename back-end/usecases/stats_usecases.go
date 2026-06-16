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
	numIndicacoes, numConversas, err := uc.repository.GetPessoas()
	if err != nil {
		return nil, err
	}
	numPessoas := numIndicacoes + numConversas

	indicacoesAprovadas, requerimentosAprovados, err := uc.repository.GetAprovados()
	if err != nil {
		return nil, err
	}

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

	numReprovados, err := uc.repository.GetReprovados()
	if err != nil {
		return nil, err
	}

	percIndicacao := 0.0
	if numPessoas > 0 {
		percIndicacao = float64(numIndicacoes) / float64(numPessoas) * 100
	}

	stat := models.Stat{
		IndicacoesAprovadas:    int(indicacoesAprovadas),
		RequerimentosAprovados: int(requerimentosAprovados),
		NumPessoas:             int(numPessoas),
		PercIndicacao:          percIndicacao,
		NumIndicacoes:          int(numIndicacoes),
		NumRequerimentos:       int(numReprovados),
		Regioes:                statsRegiao,
		Tipos:                  statsTipo,
		Categorias:             statsCategoria,
	}

	return &stat, nil
}
