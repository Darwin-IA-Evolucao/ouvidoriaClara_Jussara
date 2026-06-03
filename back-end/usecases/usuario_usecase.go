package usecases

import (
	"back-end/models"
	"back-end/repository"
)

type UsuarioUseCases struct {
	repository repository.UsuarioRepository
}

func NewUsuarioUseCases(repo repository.UsuarioRepository) UsuarioUseCases {
	return UsuarioUseCases{repository: repo}
}

func (usecase UsuarioUseCases) Login(usuario models.Usuario) (bool, error) {
	return usecase.repository.ExisteUsuario(usuario)
}
