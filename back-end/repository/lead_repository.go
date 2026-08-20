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

func (repo LeadRepository) GetAllContatosUnificados(limit, offset int) ([]models.ContatosUnificados, error) {
	query := `
		SELECT
			c.telefone,
			COALESCE(cl.nome, c.nome, c.telefone) AS nome,
			cl.cidade,
			cl.endereco,
			cl.bairro,
			cl.data_nascimento,
			COALESCE(cl.data_criacao::text, c.data_criacao::text) AS data_criacao,
			c.conversation_id,
			c.instance,
			c.campanha,
			c.ativo AS lead_ativo,
			(cb.idcliente IS NULL) AS darwin_ativo,
			(r.telefone IS NULL) AS is_gelado,
			(cl.telefone IS NOT NULL) AS is_cliente
		FROM contatos c
		LEFT JOIN cliente cl ON cl.telefone = c.telefone
		LEFT JOIN clientesbloqueados cb ON cb.idcliente = c.telefone
		LEFT JOIN (
			SELECT DISTINCT telefone FROM reclamacao
		) r ON r.telefone = c.telefone
		ORDER BY c.data_criacao DESC
	`
	if limit > 0 && offset >= 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	var leadsUnificados []models.ContatosUnificados
	err := repo.connection.Select(&leadsUnificados, query)
	if err != nil {
		return nil, err
	}
	return leadsUnificados, nil
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
