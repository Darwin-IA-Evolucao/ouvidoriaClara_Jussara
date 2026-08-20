package usecases

import (
	"back-end/config"
	"back-end/models"
	"back-end/repository"
	"database/sql"
	"fmt"
)

type LeadUseCases struct {
	repository repository.LeadRepository
}

func NewLeadUseCases(repo repository.LeadRepository) LeadUseCases {
	return LeadUseCases{
		repository: repo,
	}
}

func (usecase LeadUseCases) DesativarLead(telefone string) (sql.Result, error) {
	return usecase.repository.DesativarLead(telefone)
}

func (usecase LeadUseCases) AtivarLead(telefone string) (sql.Result, error) {
	return usecase.repository.AtivarLead(telefone)
}

func (usecase LeadUseCases) GetAllLeads(limit, offset int) ([]models.Contact, error) {
	return usecase.repository.GetAllLeads(limit, offset)
}

func (usecase LeadUseCases) GetAllContatosUnificados(limit, offset int) (*models.ContatosUnificadosResponse, error) {
	leadsUnificados, err := usecase.repository.GetAllContatosUnificados(limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := usecase.repository.GetCountContatos()
	if err != nil {
		return nil, err
	}
	usados, err := usecase.repository.GetCountContatosAtivos()
	if err != nil {
		return nil, err
	}
	planoTotal := config.GetPlanoAtual()
	return &models.ContatosUnificadosResponse{
		ContatosUnificados: leadsUnificados,
		Total:              total,
		Limite:             planoTotal,
		Usados:             usados,
		Ocupacao:           fmt.Sprintf("%d/%d", usados, planoTotal),
	}, nil
}

func (usecase LeadUseCases) GetCountContatosAtivos() (int, error) {
	return usecase.repository.GetCountContatosAtivos()
}

func (usecase LeadUseCases) GetCountContatos() (int, error) {
	return usecase.repository.GetCountContatos()
}
