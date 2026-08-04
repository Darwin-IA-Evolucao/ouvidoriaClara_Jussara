package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type CampanhaRepository struct {
	db *sqlx.DB
}

func NewCampanhaRepository(db *sqlx.DB) *CampanhaRepository {
	return &CampanhaRepository{db: db}
}

func (r *CampanhaRepository) GetCampanhaByPalavraChave(palavraChave string) (*models.Campanha, error) {
	const query = `SELECT id, palavra_chave FROM campanhas WHERE palavra_chave = $1`
	campanha := models.Campanha{}
	err := r.db.Get(&campanha, query, palavraChave)
	return &campanha, err
}

func (r *CampanhaRepository) CreateCampanha(campanha *models.Campanha) error {
	const query = `INSERT INTO campanhas (palavra_chave) VALUES ($1)`
	_, err := r.db.Exec(query, campanha.PalavraChave)
	return err
}

func (r *CampanhaRepository) UpdateCampanha(campanha *models.Campanha) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var antiga string
	err = tx.Get(&antiga, `SELECT palavra_chave FROM campanhas WHERE id = $1`, campanha.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE campanhas SET palavra_chave = $1 WHERE id = $2`, campanha.PalavraChave, campanha.ID)
	if err != nil {
		return err
	}

	if antiga != campanha.PalavraChave {
		_, err = tx.Exec(`UPDATE contatos SET campanha = $1 WHERE campanha = $2`, campanha.PalavraChave, antiga)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CampanhaRepository) DeleteCampanha(id int) error {
	const query = `DELETE FROM campanhas WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CampanhaRepository) GetAllCampanhas() ([]models.Campanha, error) {
	const query = `
		SELECT c.id, c.palavra_chave, COUNT(ct.telefone) AS qtd_contatos
		FROM campanhas c
		LEFT JOIN contatos ct ON ct.campanha = c.palavra_chave
		GROUP BY c.id, c.palavra_chave
		ORDER BY c.id`
	campanhas := []models.Campanha{}
	err := r.db.Select(&campanhas, query)
	return campanhas, err
}

func (r *CampanhaRepository) GetCampanhaByID(id int) (*models.Campanha, error) {
	const query = `SELECT id, palavra_chave FROM campanhas WHERE id = $1`
	campanha := models.Campanha{}
	err := r.db.Get(&campanha, query, id)
	return &campanha, err
}

func (r *CampanhaRepository) GetContatosByCampanha(palavraChave string) ([]models.ContatoCampanha, error) {
	const query = `SELECT telefone, COALESCE(nome, '') AS nome FROM contatos WHERE campanha = $1 ORDER BY data_criacao DESC`
	contatos := []models.ContatoCampanha{}
	err := r.db.Select(&contatos, query, palavraChave)
	return contatos, err
}

func (r *CampanhaRepository) ExistsCampanha(palavraChave string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM campanhas WHERE palavra_chave = $1)`
	var exists bool
	err := r.db.Get(&exists, query, palavraChave)
	return exists, err
}
