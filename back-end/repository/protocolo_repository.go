package repository

import (
	"github.com/jmoiron/sqlx"
)

type ProtocoloRepository struct {
	connection *sqlx.DB
}

func NewProtocoloRepository(conn *sqlx.DB) ProtocoloRepository {
	return ProtocoloRepository{connection: conn}
}

type DadosReclamacao struct {
	Nome      string `db:"nome"`
	Categoria string `db:"categoria"`
	Telefone  string `db:"telefone"`
}

func (repo ProtocoloRepository) CreateProtocolo(numero, idReclamacao int64) error {
	const query = `INSERT INTO protocolo (numero, idreclamacao) VALUES ($1, $2)`
	_, err := repo.connection.Exec(query, numero, idReclamacao)
	return err
}

func (repo ProtocoloRepository) MarcarResolvido(idReclamacao int64) error {
	const query = `UPDATE reclamacao SET resolvido = TRUE WHERE id = $1`
	_, err := repo.connection.Exec(query, idReclamacao)
	return err
}

func (repo ProtocoloRepository) GetDadosReclamacao(idReclamacao int64) (*DadosReclamacao, error) {
	const query = `SELECT nome, categoria, telefone FROM reclamacao WHERE id = $1`
	var dados DadosReclamacao
	err := repo.connection.Get(&dados, query, idReclamacao)
	if err != nil {
		return nil, err
	}
	return &dados, nil
}

func (repo ProtocoloRepository) MarcarAvisado(idReclamacao int64) error {
	const query = `UPDATE protocolo SET avisado = TRUE WHERE idreclamacao = $1`
	_, err := repo.connection.Exec(query, idReclamacao)
	return err
}
