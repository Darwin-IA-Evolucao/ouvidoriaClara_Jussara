package repository

import (
	"back-end/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UsuarioRepository struct {
	connection *sqlx.DB
}

func NewUsuarioRepository(conn *sqlx.DB) UsuarioRepository {
	return UsuarioRepository{
		connection: conn,
	}
}

func (repo UsuarioRepository) GetUsuarioLogin(usuario models.Usuario) (int, string, error) {
	const query = `SELECT id, role 
	          FROM usuarios 
	          WHERE (celular = $1 AND senha = $2 AND ativo = TRUE) 
	             OR (role = 'root' AND celular = $1 AND senha = $2)`

	var id int
	var role string
	err := repo.connection.QueryRowx(query, usuario.Celular, usuario.Senha).Scan(&id, &role)
	if err != nil {
		return 0, "", fmt.Errorf("CREDENCIAIS INVALIDAS")
	}

	return id, role, nil
}
