package usecases

import (
	"back-end/models"
	"back-end/repository"
	"database/sql"
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

func (usecase LeadUseCases) GetCountContatosAtivos() (int, error) {
	return usecase.repository.GetCountContatosAtivos()
}

func (usecase LeadUseCases) GetCountContatos() (int, error) {
	return usecase.repository.GetCountContatos()
}
