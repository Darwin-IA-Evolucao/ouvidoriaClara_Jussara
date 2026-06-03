package repository

import (
	"github.com/jmoiron/sqlx"
)

type StatsRepository struct {
	connection *sqlx.DB
}

func NewStatsRepository(conn *sqlx.DB) StatsRepository {
	return StatsRepository{connection: conn}
}

func (repo StatsRepository) GetPessoas() (int64, int64, error) {
	const query = `SELECT
		(SELECT COUNT(DISTINCT telefone) FROM reclamacao) AS pessoas_com_reclamacao,
		(SELECT COUNT(DISTINCT telefone) FROM contatos WHERE telefone NOT IN (SELECT telefone FROM reclamacao)) AS pessoas_sem_reclamacao;`
	var comReclamacao, semReclamacao int64
	err := repo.connection.QueryRow(query).Scan(&comReclamacao, &semReclamacao)
	return comReclamacao, semReclamacao, err
}

func (repo StatsRepository) GetAprovadosEProtocoladas() (int64, int64, int64, int64, error) {
	const query = `SELECT
		(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'indicação') AS indicacoes_aprovadas,
		(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'requerimento') AS requerimentos_aprovados,
		(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'ofício') AS oficios_aprovados,
		(SELECT COUNT(idReclamacao) FROM Protocolo) AS protocoladas;`
	var indAprovadas, reqAprovados, ofAprovados, protocoladas int64
	err := repo.connection.QueryRow(query).Scan(&indAprovadas, &reqAprovados, &ofAprovados, &protocoladas)
	return indAprovadas, reqAprovados, ofAprovados, protocoladas, err
}

func (repo StatsRepository) GetCountByCategoriaRegiao(categoria string, regiao int) (int, error) {
	const query = `SELECT COUNT(id) FROM Reclamacao WHERE categoria = $1 AND Regiao = $2`
	var count int
	err := repo.connection.Get(&count, query, categoria, regiao)
	return count, err
}

func (repo StatsRepository) GetTipos() (int64, int64, int64, error) {
	const query = `SELECT
		(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'indicação') AS numIndicacao,
		(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'ofício') AS numOficios,
		(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'requerimento') AS numRequerimentos;`
	var indicacoes, numOficios, numRequerimentos int64
	err := repo.connection.QueryRow(query).Scan(&indicacoes, &numOficios, &numRequerimentos)
	return indicacoes, numOficios, numRequerimentos, err
}

func (repo StatsRepository) GetReprovados() (int64, error) {
	const query = `SELECT COUNT(id) FROM Reclamacao WHERE status = 'reprovado';`
	var reprovados int64
	err := repo.connection.Get(&reprovados, query)
	return reprovados, err
}
