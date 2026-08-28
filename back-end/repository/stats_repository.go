package repository

import (
	"back-end/models"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type StatsRepository struct {
	connection *sqlx.DB
}

func NewStatsRepository(conn *sqlx.DB) StatsRepository {
	return StatsRepository{connection: conn}
}

func whereDataCriacao(inicio, fim *time.Time) (string, []any) {
	conds := []string{}
	args := []any{}
	if inicio != nil {
		args = append(args, *inicio)
		conds = append(conds, fmt.Sprintf("data_criacao >= $%d", len(args)))
	}
	if fim != nil {
		args = append(args, *fim)
		conds = append(conds, fmt.Sprintf("data_criacao < $%d", len(args)))
	}
	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func (repo StatsRepository) GetPessoas(inicio, fim *time.Time) (int64, int64, error) {
	where, args := whereDataCriacao(inicio, fim)
	contatosWhere := where
	notIn := "telefone NOT IN (SELECT telefone FROM reclamacao" + where + ")"
	if contatosWhere == "" {
		contatosWhere = " WHERE " + notIn
	} else {
		contatosWhere += " AND " + notIn
	}
	query := fmt.Sprintf(`SELECT
		(SELECT COUNT(DISTINCT telefone) FROM reclamacao%s) AS pessoas_com_reclamacao,
		(SELECT COUNT(DISTINCT telefone) FROM contatos%s) AS pessoas_sem_reclamacao;`, where, contatosWhere)
	var comReclamacao, semReclamacao int64
	err := repo.connection.QueryRow(query, args...).Scan(&comReclamacao, &semReclamacao)
	return comReclamacao, semReclamacao, err
}

func (repo StatsRepository) GetCountByCategoria(inicio, fim *time.Time) ([]models.StatsCategoria, error) {
	where, args := whereDataCriacao(inicio, fim)
	query := `SELECT categoria, count(*) AS qtd_categoria FROM reclamacao` + where + ` GROUP BY categoria;`
	var stats []models.StatsCategoria
	err := repo.connection.Select(&stats, query, args...)
	return stats, err
}

func (repo StatsRepository) GetCountByRegiao(inicio, fim *time.Time) ([]models.StatsRegiao, error) {
	where, args := whereDataCriacao(inicio, fim)
	query := ` SELECT
						COALESCE(detalhes->>'regiao', 'Sem Regiao Definida') AS regiao,
						COUNT(*) AS qtd_regiao
					FROM reclamacao` + where + `
					GROUP BY COALESCE(detalhes->>'regiao', 'Sem Regiao Definida')
					ORDER BY regiao;`
	var stats []models.StatsRegiao
	err := repo.connection.Select(&stats, query, args...)
	return stats, err
}

func (repo StatsRepository) GetCountByTipo(inicio, fim *time.Time) ([]models.StatsTipo, error) {
	where, args := whereDataCriacao(inicio, fim)
	query := `SELECT
					COALESCE(NULLIF(tipo,''), 'Sem Tipo Definido') AS tipo,
					COUNT(*) as qtd_tipo
					FROM reclamacao` + where + `
					GROUP BY COALESCE(NULLIF(tipo,''), 'Sem Tipo Definido')
					ORDER BY tipo;`
	var stats []models.StatsTipo
	err := repo.connection.Select(&stats, query, args...)
	return stats, err
}

func (repo StatsRepository) GetCountByTipoAndStatus(inicio, fim *time.Time) (models.StatsByTipoAndStatus, error) {
	where, args := whereDataCriacao(inicio, fim)
	query := `
		SELECT
		COUNT(*) FILTER (
			WHERE tipo ILIKE 'indicacao'
			AND status ILIKE 'aprovado'
		) AS indicacoes_aprovadas,

		COUNT(*) FILTER (
			WHERE tipo ILIKE 'indicacao'
		) AS total_indicacoes,

		COUNT(*) FILTER (
			WHERE tipo ILIKE 'requerimento'
			AND status ILIKE 'aprovado'
		) AS requerimentos_aprovados,

		COUNT(*) FILTER (
			WHERE tipo ILIKE 'requerimento'
		) AS total_requerimentos,

		COUNT(*) AS total_reclamacao
		FROM reclamacao` + where + `;`
	var stats models.StatsByTipoAndStatus
	err := repo.connection.Get(&stats, query, args...)
	return stats, err
}
