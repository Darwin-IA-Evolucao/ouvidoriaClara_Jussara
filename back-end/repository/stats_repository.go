package repository

import (
	"back-end/models"

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

func (repo StatsRepository) GetCountByCategoria() ([]models.StatsCategoria, error) {
	const query = `SELECT categoria, count(*) AS qtd_categoria FROM reclamacao GROUP BY categoria;`
	var stats []models.StatsCategoria
	err := repo.connection.Select(&stats, query)
	return stats, err
}

func (repo StatsRepository) GetCountByRegiao() ([]models.StatsRegiao, error) {
	const query = ` SELECT
						COALESCE(detalhes->>'regiao', 'Sem Regiao Definida') AS regiao,
						COUNT(*) AS qtd_regiao
					FROM reclamacao
					GROUP BY COALESCE(detalhes->>'regiao', 'Sem Regiao Definida')
					ORDER BY regiao;`
	var stats []models.StatsRegiao
	err := repo.connection.Select(&stats, query)
	return stats, err
}

func (repo StatsRepository) GetCountByTipo() ([]models.StatsTipo, error) {
	const query = `SELECT
					COALESCE(NULLIF(tipo,''), 'Sem Tipo Definido') AS tipo,
					COUNT(*) as qtd_tipo
					FROM reclamacao
					GROUP BY COALESCE(NULLIF(tipo,''), 'Sem Tipo Definido')
					ORDER BY tipo;`
	var stats []models.StatsTipo
	err := repo.connection.Select(&stats, query)
	return stats, err
}

func (repo StatsRepository) GetCountByTipoAndStatus() (models.StatsByTipoAndStatus, error) {
	const query = `
		SELECT
		COUNT(*) FILTER (
			WHERE tipo = 'indicacao'
			AND status = 'aprovado'
		) AS indicacoes_aprovadas,

		COUNT(*) FILTER (
			WHERE tipo = 'indicacao'
		) AS total_indicacoes,

		COUNT(*) FILTER (
			WHERE tipo = 'requerimento'
			AND status = 'Aprovado'
		) AS requerimentos_aprovados,

		COUNT(*) FILTER (
			WHERE tipo = 'requerimento'
		) AS total_requerimentos,

		COUNT(*) AS total_reclamacao
		FROM reclamacao;
	`
	var stats models.StatsByTipoAndStatus
	err := repo.connection.Get(&stats, query)
	return stats, err
}
