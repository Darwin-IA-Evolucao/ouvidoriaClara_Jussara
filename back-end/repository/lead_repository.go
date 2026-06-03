package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type LeadRepository struct {
	connection *sqlx.DB
}

func NewLeadRepository(conn *sqlx.DB) LeadRepository {
	return LeadRepository{connection: conn}
}

func (repo LeadRepository) GetAllLeads() ([]models.Lead, error) {
	const query = `SELECT nome, telefone, ativo FROM contatos`
	var leads []models.Lead
	err := repo.connection.Select(&leads, query)
	if err != nil {
		return nil, err
	}
	return leads, nil
}

func (repo LeadRepository) DesativarLead(telefone string) error {
	const query = `UPDATE contatos SET ativo = false WHERE telefone = $1`
	_, err := repo.connection.Exec(query, telefone)
	return err
}
