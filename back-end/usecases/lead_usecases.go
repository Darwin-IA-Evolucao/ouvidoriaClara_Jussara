package usecases

import (
	"back-end/config"
	"back-end/models"
	"back-end/repository"
	"strconv"
)

type LeadUseCases struct {
	repository repository.LeadRepository
}

func NewLeadUseCases(repo repository.LeadRepository) LeadUseCases {
	return LeadUseCases{repository: repo}
}

func (uc LeadUseCases) GetLeads() ([]models.Lead, error) {
	return uc.repository.GetAllLeads()
}

func (uc LeadUseCases) GetAllLeads() ([]models.Lead, string, error) {
	leads, err := uc.repository.GetAllLeads()
	if err != nil {
		return nil, "", err
	}
	count := 0
	for _, lead := range leads {
		if lead.Ativo {
			count++
		}
	}
	ocupacao := strconv.Itoa(count) + "/" + strconv.Itoa(config.GetPlanoLimite())
	return leads, ocupacao, nil
}

func (uc LeadUseCases) DesativarLead(telefone string) error {
	return uc.repository.DesativarLead(telefone)
}
