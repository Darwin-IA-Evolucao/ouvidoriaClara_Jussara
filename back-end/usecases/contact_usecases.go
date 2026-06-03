package usecases

import (
	"back-end/models"
	"back-end/repository"
)

type ContatoUseCase struct {
	repo repository.ContatoRepo
}

func NewContatoUseCase(repo repository.ContatoRepo) *ContatoUseCase {
	return &ContatoUseCase{repo: repo}
}

func (u *ContatoUseCase) GetAllContatos() ([]models.Contact, error) {
	return u.repo.GetAllContatos()
}

func (u *ContatoUseCase) GetContatoByTelefone(telefone string) (*models.Contato, error) {
	return u.repo.GetContatoByTelefone(telefone)
}

func (u *ContatoUseCase) SetConversationId(telefone string, conversationId string) error {
	return u.repo.SetConversationId(telefone, conversationId)
}

func (u *ContatoUseCase) GetConversationIdByTelefone(telefone string) (*models.GetConversationIdResponse, error) {
	contato, err := u.repo.GetContatoByTelefone(telefone)
	if err != nil {
		return nil, err
	}

	var response models.GetConversationIdResponse
	response.Exists = contato.ConversationId != nil

	if contato.ConversationId != nil {
		response.ConversationId = *contato.ConversationId
	} else {
		response.ConversationId = ""
	}

	return &response, nil
}
