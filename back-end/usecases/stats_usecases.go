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

	indicacoesAprovadas, requerimentosAprovados, oficiosAprovados, protocoladas, err := uc.repository.GetAprovadosEProtocoladas()
	if err != nil {
		return nil, err
	}

	var regioes []models.Regiao
	for i := 1; i <= 6; i++ {
		mapeamento := make(map[string]int)
		for _, cat := range categorias {
			contador, err := uc.repository.GetCountByCategoriaRegiao(cat, i)
			if err != nil {
				return nil, err
			}
			mapeamento[cat] = contador
		}
		regioes = append(regioes, models.Regiao{ID: i, Distribuicao: mapeamento})
	}

	indicacoes, numOficios, numRequerimentos, err := uc.repository.GetTipos()
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

	percProtocolacao := 0.0
	denominator := indicacoesAprovadas + requerimentosAprovados
	if denominator > 0 {
		percProtocolacao = float64(protocoladas) / float64(denominator) * 100
	}

	stat := models.Stat{
		IndicacoesAprovadas:    int(indicacoesAprovadas),
		RequerimentosAprovados: int(requerimentosAprovados),
		OficiosAprovados:       int(oficiosAprovados),
		NumPessoas:             int(numPessoas),
		Indicacoes:             int(indicacoes),
		NumRequerimentos:       int(numRequerimentos),
		NumOficios:             int(numOficios),
		PercIndicacao:          percIndicacao,
		NumReclamacoes:         int(indicacoes + numRequerimentos),
		NumIndicacoes:          int(numIndicacoes),
		NumReprovados:          int(numReprovados),
		Protocoladas:           int(protocoladas),
		PercProtocolacao:       percProtocolacao,
		Regioes:                regioes,
	}

	return &stat, nil
}
