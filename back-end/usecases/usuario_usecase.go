package usecases

import (
	"back-end/models"
	"back-end/repository"
)

type UsuarioUseCases struct {
	repository repository.UsuarioRepository
}

func NewUsuarioUseCases(repo repository.UsuarioRepository) UsuarioUseCases {
	return UsuarioUseCases{
		repository: repo,
	}
}

func (usecase UsuarioUseCases) GetUsuarioLogin(usuario models.Usuario) (int, string, bool, error) {
	id, role, err := usecase.repository.GetUsuarioLogin(usuario)
	if err != nil {
		return 0, "", false, err
	}
	ehRoot := role == "root"
	return id, role, ehRoot, nil
}
