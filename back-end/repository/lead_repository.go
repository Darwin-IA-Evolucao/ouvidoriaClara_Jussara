package repository

import (
	"back-end/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type LeadRepository struct {
	connection *sqlx.DB
}

func NewLeadRepository(conn *sqlx.DB) LeadRepository {
	return LeadRepository{
		connection: conn,
	}
}

func (repo LeadRepository) DesativarLead(telefone string) (sql.Result, error) {
	const query = `UPDATE contatos SET ativo = false WHERE telefone = $1`

	res, err := repo.connection.Exec(query, telefone)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo LeadRepository) AtivarLead(telefone string) (sql.Result, error) {
	const query = `UPDATE contatos SET ativo = true WHERE telefone = $1`

	res, err := repo.connection.Exec(query, telefone)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo LeadRepository) GetAllLeads(limit, offset int) ([]models.Contact, error) {
	query := `SELECT * FROM contatos`
	if limit > 0 && offset >= 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	var leads []models.Contact
	err := repo.connection.Select(&leads, query)
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (repo LeadRepository) GetCountContatosAtivos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos WHERE ativo = true`

	var count int
	err := repo.connection.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo LeadRepository) GetCountContatos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos`

	var count int
	err := repo.connection.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
