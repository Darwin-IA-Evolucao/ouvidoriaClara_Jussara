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

func (repo StatsRepository) GetAprovados() (int64, int64, error) {
	const query = `SELECT
		(SELECT COUNT(*) FROM reclamacao WHERE status = 'aprovado' AND tipo = 'indicação') AS indicacoes_aprovadas,
		(SELECT COUNT(*) FROM reclamacao WHERE status = 'aprovado' AND tipo = 'requerimento') AS requerimentos_aprovados;`
	var indAprovadas, reqAprovados, ofAprovados, protocoladas int64
	err := repo.connection.QueryRow(query).Scan(&indAprovadas, &reqAprovados, &ofAprovados, &protocoladas)
	return indAprovadas, reqAprovados, err
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

func (repo StatsRepository) GetReprovados() (int64, error) {
	const query = `SELECT COUNT(*) FROM reclamacao WHERE status = 'reprovado';`
	var reprovados int64
	err := repo.connection.Get(&reprovados, query)
	return reprovados, err
}
