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

func (r *ContatoRepo) GetAllContatos()([]models.Contact, error){
	const query = `SELECT * FROM contatos`
	var contatos []models.Contact
	err := r.db.Select(&contatos, query)
	return contatos, err
}

func (r *ContatoRepo)	GetContatoByTelefone(telefone string) (*models.Contato, error) {
	const query = `SELECT * FROM contatos WHERE telefone = $1`
	var contato models.Contato
	err := r.db.Get(&contato, query, telefone)
	if err != nil {
		return nil, err
	}
	return &contato, nil
}
func (r ContatoRepo) GetClienteBloqueadoById(telefoneCliente string) error{
	query := `SELECT idcliente FROM clientesbloqueados WHERE idcliente = $1`

	var bloqueado string
	err := r.db.Get(&bloqueado, query, telefoneCliente)
	if err != nil {
		return err
	}
	return nil
}

func (r ContatoRepo) GetCountContatos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos`
	var count int
	err := r.db.Get(&count, query)
	return count, err
}

func (r *ContatoRepo) UpdateInstanciaContato(telefone string, instancia string) error {
	const query = `UPDATE contatos SET instance = $2 WHERE telefone = $1`
	_, err := r.db.Exec(query, telefone, instancia)
	return err
}

func (r *ContatoRepo) CreateContato(contato *models.Contato) error {
	const query = `INSERT INTO contatos (telefone, nome, conversation_id, instance) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, contato.Telefone, contato.Nome, contato.ConversationId, contato.Instance)
	return err
}

func (r *ContatoRepo) SetConversationId(telefone string, conversationId string) error {
	const query = `UPDATE contatos SET conversation_id = $2 WHERE telefone = $1`
	_, err := r.db.Exec(query, telefone, conversationId)
	return err
}