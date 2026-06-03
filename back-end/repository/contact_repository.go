package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type ContatoRepo struct {
	db *sqlx.DB
}

func NewContatoRepo(db *sqlx.DB) *ContatoRepo {
	return &ContatoRepo{db: db}
}

func (r *ContatoRepo) GetCountContatosAtivos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos WHERE ativo = true`
	var count int
	err := r.db.Get(&count, query)
	return count, err
}

func (r *ContatoRepo) GetContatoByTelefone(telefone string) (*models.Contact, error) {
	const query = `SELECT nome, telefone, conversation_id FROM contatos WHERE telefone = $1`
	var contato models.Contact
	err := r.db.QueryRowx(query, telefone).Scan(&contato.Nome, &contato.Telefone, &contato.ConversationID)
	if err != nil {
		return nil, err
	}
	return &contato, nil
}

func (r *ContatoRepo) CreateContato(nome, telefone string) error {
	const query = `INSERT INTO contatos (nome, telefone, conversation_id, ativo) VALUES ($1, $2, $3, true)`
	_, err := r.db.Exec(query, nome, telefone, nil)
	return err
}

func (r *ContatoRepo) SetConversationId(telefone, conversationId string) error {
	const query = `UPDATE contatos SET conversation_id = $1 WHERE telefone = $2`
	_, err := r.db.Exec(query, conversationId, telefone)
	return err
}
