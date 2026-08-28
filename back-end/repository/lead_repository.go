package repository

import (
	"back-end/models"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

var nonDigit = regexp.MustCompile(`\D`)

const contatosUnificadosFrom = `
		FROM contatos c
		LEFT JOIN cliente cl ON cl.telefone = c.telefone
		LEFT JOIN clientesbloqueados cb ON cb.idcliente = c.telefone
		LEFT JOIN (
			SELECT DISTINCT telefone FROM reclamacao
		) r ON r.telefone = c.telefone`

type LeadRepository struct {
	connection *sqlx.DB
}

func NewLeadRepository(conn *sqlx.DB) LeadRepository {
	return LeadRepository{
		connection: conn,
	}
}

func (repo LeadRepository) DesativarLead(telefone string) (sql.Result, error) {
	const query = `UPDATE contatos SET ativo = false WHERE telefone = $1`

	res, err := repo.connection.Exec(query, telefone)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo LeadRepository) AtivarLead(telefone string) (sql.Result, error) {
	const query = `UPDATE contatos SET ativo = true WHERE telefone = $1`

	res, err := repo.connection.Exec(query, telefone)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo LeadRepository) GetAllLeads(limit, offset int) ([]models.Contact, error) {
	query := `SELECT * FROM contatos`
	if limit > 0 && offset >= 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
	}
	var leads []models.Contact
	err := repo.connection.Select(&leads, query)
	if err != nil {
		return nil, err
	}

	return leads, nil
}

func (repo LeadRepository) filtroContatosUnificados(f models.ContatosUnificadosFiltro) (string, []any) {
	conds := []string{}
	args := []any{}
	if f.Search != "" {
		args = append(args, "%"+f.Search+"%")
		nameArg := len(args)
		digits := nonDigit.ReplaceAllString(f.Search, "")
		if digits != "" {
			args = append(args, "%"+digits+"%")
			conds = append(conds, fmt.Sprintf(`(COALESCE(cl.nome, c.nome, c.telefone) ILIKE $%d OR c.telefone LIKE $%d)`, nameArg, len(args)))
		} else {
			conds = append(conds, fmt.Sprintf(`COALESCE(cl.nome, c.nome, c.telefone) ILIKE $%d`, nameArg))
		}
	}
	if f.Darwin != nil {
		if *f.Darwin {
			conds = append(conds, "cb.idcliente IS NULL")
		} else {
			conds = append(conds, "cb.idcliente IS NOT NULL")
		}
	}
	if f.LeadAtivo != nil {
		args = append(args, *f.LeadAtivo)
		conds = append(conds, fmt.Sprintf("c.ativo = $%d", len(args)))
	}
	if f.Gelo != nil {
		if *f.Gelo {
			conds = append(conds, "r.telefone IS NULL")
		} else {
			conds = append(conds, "r.telefone IS NOT NULL")
		}
	}
	if f.Reclamacao != nil {
		if *f.Reclamacao {
			conds = append(conds, "r.telefone IS NOT NULL")
		} else {
			conds = append(conds, "r.telefone IS NULL")
		}
	}
	if f.Inicio != nil {
		args = append(args, *f.Inicio)
		conds = append(conds, fmt.Sprintf("COALESCE(cl.data_criacao, c.data_criacao::timestamp) >= $%d", len(args)))
	}
	if f.Fim != nil {
		args = append(args, *f.Fim)
		conds = append(conds, fmt.Sprintf("COALESCE(cl.data_criacao, c.data_criacao::timestamp) < $%d", len(args)))
	}
	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func (repo LeadRepository) GetAllContatosUnificados(limit, offset int, filtro models.ContatosUnificadosFiltro) ([]models.ContatosUnificados, error) {
	where, args := repo.filtroContatosUnificados(filtro)
	query := `
		SELECT
			c.telefone,
			COALESCE(cl.nome, c.nome, c.telefone) AS nome,
			cl.cidade,
			cl.endereco,
			cl.bairro,
			cl.data_nascimento,
			COALESCE(cl.data_criacao::text, c.data_criacao::text) AS data_criacao,
			c.conversation_id,
			c.instance,
			c.campanha,
			c.ativo AS lead_ativo,
			(cb.idcliente IS NULL) AS darwin_ativo,
			(r.telefone IS NULL) AS is_gelado,
			(cl.telefone IS NOT NULL) AS is_cliente,
			(r.telefone IS NOT NULL) AS has_reclamacao
		` + contatosUnificadosFrom + where + `
		ORDER BY c.data_criacao DESC, c.telefone ASC`
	if limit > 0 && offset >= 0 {
		args = append(args, limit, offset)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}
	var leadsUnificados []models.ContatosUnificados
	err := repo.connection.Select(&leadsUnificados, query, args...)
	if err != nil {
		return nil, err
	}
	if leadsUnificados == nil {
		leadsUnificados = []models.ContatosUnificados{}
	}
	return leadsUnificados, nil
}

func (repo LeadRepository) GetCountContatosUnificados(filtro models.ContatosUnificadosFiltro) (int, error) {
	where, args := repo.filtroContatosUnificados(filtro)
	query := `SELECT COUNT(*) ` + contatosUnificadosFrom + where
	var count int
	err := repo.connection.Get(&count, query, args...)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo LeadRepository) GetCountContatosAtivos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos WHERE ativo = true`

	var count int
	err := repo.connection.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo LeadRepository) GetCountContatos() (int, error) {
	const query = `SELECT COUNT(*) FROM contatos`

	var count int
	err := repo.connection.Get(&count, query)
	if err != nil {
		return 0, err
	}
	return count, nil
}
