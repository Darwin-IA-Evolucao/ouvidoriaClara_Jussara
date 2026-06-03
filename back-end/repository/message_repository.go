package repository

import (
	"back-end/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type MensagemRepo struct {
	db *sqlx.DB
}

func NewMensagemRepo(db *sqlx.DB) *MensagemRepo {
	return &MensagemRepo{db: db}
}

func (r *MensagemRepo) GetCountContatosAtivos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos WHERE ativo = true`
	var count int
	err := r.db.Get(&count, query)
	return count, err
}

func (r *MensagemRepo) GetCountContatoAtivoByTelefone(telefone string) (int, error) {
	const query = `SELECT COUNT(*) FROM contatos WHERE telefone = $1 AND ativo = true`
	var count int
	err := r.db.Get(&count, query, telefone)
	return count, err
}

func (r *MensagemRepo) CreateMensagem(telefone, conteudo string) error {
	const query = `INSERT INTO mensagens (telefone, conteudo) VALUES ($1, $2)`
	_, err := r.db.Exec(query, telefone, conteudo)
	return err
}

func (r *MensagemRepo) GetCountConversaHoje(telefone, data string) (int, error) {
	const query = `SELECT COUNT(*) FROM conversas WHERE telefone = $1 AND data = $2`
	var count int
	err := r.db.Get(&count, query, telefone, data)
	return count, err
}

func (r *MensagemRepo) CreateConversa(telefone, data string) error {
	const query = `INSERT INTO conversas (telefone, data) VALUES ($1, $2)`
	_, err := r.db.Exec(query, telefone, data)
	return err
}

func (r *MensagemRepo) AtivarContato(telefone string) error {
	const query = `UPDATE contatos SET ativo = true WHERE telefone = $1`
	_, err := r.db.Exec(query, telefone)
	return err
}

func (r *MensagemRepo) GetCountMensagens(telefone string) (int, error) {
	const query = `SELECT COUNT(*) FROM mensagens WHERE telefone = $1`
	var count int
	err := r.db.Get(&count, query, telefone)
	return count, err
}

func (r *MensagemRepo) GetMensagensByTelefone(telefone string) ([]models.Message, error) {
	const query = `SELECT telefone, conteudo FROM mensagens WHERE telefone = $1`
	var mensagens []models.Message
	err := r.db.Select(&mensagens, query, telefone)
	return mensagens, err
}

func (r *MensagemRepo) DeleteMensagens(telefone string) error {
	const query = `DELETE FROM mensagens WHERE telefone = $1`
	_, err := r.db.Exec(query, telefone)
	return err
}

func (r *MensagemRepo) GetConversationId(telefone string) (string, error) {
	const query = `SELECT conversation_id FROM contatos WHERE telefone = $1`
	var conversationID sql.NullString
	err := r.db.Get(&conversationID, query, telefone)
	if err != nil {
		return "", err
	}
	if conversationID.Valid {
		return conversationID.String, nil
	}
	return "", nil
}
