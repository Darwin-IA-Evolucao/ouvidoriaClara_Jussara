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
	const query = `SELECT * FROM campanhas WHERE palavra_chave = $1`
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
	const query = `UPDATE campanhas SET palavra_chave = $1 WHERE id = $2`
	_, err := r.db.Exec(query, campanha.PalavraChave, campanha.ID)
	return err
}

func (r *CampanhaRepository) DeleteCampanha(id int) error {
	const query = `DELETE FROM campanhas WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *CampanhaRepository) GetAllCampanhas() ([]models.Campanha, error) {
	const query = `SELECT * FROM campanhas`
	campanhas := []models.Campanha{}
	err := r.db.Select(&campanhas, query)
	return campanhas, err
}

func (r *CampanhaRepository) ExistsCampanha(palavraChave string) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM campanhas WHERE palavra_chave = $1)`
	var exists bool
	err := r.db.Get(&exists, query, palavraChave)
	return exists, err
}
