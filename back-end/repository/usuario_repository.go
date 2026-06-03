package repository

import (
	"back-end/models"

	"github.com/jmoiron/sqlx"
)

type UsuarioRepository struct {
	connection *sqlx.DB
}

func NewUsuarioRepository(conn *sqlx.DB) UsuarioRepository {
	return UsuarioRepository{connection: conn}
}

func (repo UsuarioRepository) ExisteUsuario(usuario models.Usuario) (bool, error) {
	const query = `SELECT COUNT(*) FROM usuarios WHERE usuario = $1 AND senha = $2`
	var count int
	err := repo.connection.Get(&count, query, usuario.Usuario, usuario.Senha)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
