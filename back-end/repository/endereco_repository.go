package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type EnderecoRepository struct {
	connection *sqlx.DB
}

func NewEnderecoRepository(conn *sqlx.DB) EnderecoRepository {
	return EnderecoRepository{connection: conn}
}

func (repo EnderecoRepository) GetRegiaoByLogradouro(input string) (int, error) {
	const query = `
		SELECT regiao FROM enderecos
		WHERE logradouro ILIKE '%' || $1 || '%'
		ORDER BY LENGTH(logradouro) ASC LIMIT 1;`
	var regiao int
	err := repo.connection.Get(&regiao, query, input)
	return regiao, err
}

func (repo EnderecoRepository) GetAllLogradouros() ([]models.Logradouro, error) {
	const query = `SELECT logradouro, bairro, regiao FROM enderecos ORDER BY logradouro;`
	var logradouros []models.Logradouro
	err := repo.connection.Select(&logradouros, query)
	return logradouros, err
}

func (repo EnderecoRepository) CreateEndereco(endereco models.Endereco) error {
	const query = `INSERT INTO enderecos (logradouro, bairro, regiao) VALUES ($1, $2, $3)`
	_, err := repo.connection.Exec(query, endereco.Logradouro, endereco.Bairro, endereco.Regiao)
	return err
}
