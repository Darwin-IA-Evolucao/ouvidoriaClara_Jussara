package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type MensagemRepo struct {
	db *sqlx.DB
}

func NewMensagemRepo(db *sqlx.DB) *MensagemRepo {
	return &MensagemRepo{db: db}
}

// -------- CONTATO --------
func (r *MensagemRepo) GetContatoByTelefone(telefone string) (*models.Contato, error) {
	return NewContatoRepo(r.db).GetContatoByTelefone(telefone)
}
func (r *MensagemRepo) UpdateInstanciaContato(telefone string, instancia string) error {
	return NewContatoRepo(r.db).UpdateInstanciaContato(telefone, instancia)
}

func (r *MensagemRepo) CreateContato(contato *models.Contato) error {
	return NewContatoRepo(r.db).CreateContato(contato)
}

func (r *MensagemRepo) GetClienteBloqueadoById(telefoneCliente string) error {
	return NewContatoRepo(r.db).GetClienteBloqueadoById(telefoneCliente)
}

func (r *MensagemRepo) GetCountContatos() (int, error) {
	return NewContatoRepo(r.db).GetCountContatos()
}

// -------- ATIVIDADE CLIENTE --------
func (r *MensagemRepo) UpdateUltimaInteracao(telefone string) error {
	const query = `INSERT INTO atividade_clientes (telefone) VALUES ($1) ON CONFLICT (telefone) DO UPDATE SET ultima_interacao = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query, telefone)
	return err
}

// -------- MENSAGEM --------
func (r *MensagemRepo) CreateMensagem(mensagem *models.Mensagem) error {
	const query = `INSERT INTO mensagens (telefone, conteudo) VALUES ($1, $2) RETURNING id`
	err := r.db.Get(&mensagem.ID, query, mensagem.Telefone, mensagem.Conteudo)
	return err
}
func (r *MensagemRepo) GetMensagensNaoEnviadasBytelefone(telefone string) ([]models.Mensagem, error) {
	const query = `SELECT * FROM mensagens WHERE telefone = $1 AND foienviado = false`
	var mensagens []models.Mensagem
	err := r.db.Select(&mensagens, query, telefone)
	return mensagens, err
}
func (r *MensagemRepo) SetFoiEnviado(id int) error {
	const query = `UPDATE mensagens SET foienviado = true WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
func (r *MensagemRepo) DeleteMensagens(telefone string) error {
	const query = `DELETE FROM mensagens WHERE telefone = $1 AND foienviado = true`
	_, err := r.db.Exec(query, telefone)
	return err
}
func (r *MensagemRepo) GetCountMensagens(telefone string) (int, error) {
	const query = `SELECT COUNT(*) FROM mensagens WHERE telefone = $1 AND foienviado = false`
	var count int
	err := r.db.Get(&count, query, telefone)
	return count, err
}

// -------- BLOQUEIO --------
func (r *MensagemRepo) GetAvisoPlano() (models.AvisoPlanoAtigido, error) {
	const query = `SELECT * FROM aviso_plano_atingido;`
	var aviso models.AvisoPlanoAtigido
	err := r.db.Get(&aviso, query)
	return aviso, err
}

func (r *MensagemRepo) SetAvisado() error {
	const query = `UPDATE aviso_plano_atingido SET avisado = true, data = CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query)
	return err
}
