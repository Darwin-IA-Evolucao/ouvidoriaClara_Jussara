package usecases

import (
	"back-end/apperror"
	"back-end/config"
	"back-end/models"
	"back-end/repository"
	"database/sql"
	"errors"
)

type ContatoUseCase struct {
	repo repository.ContatoRepo
}

func NewContatoUseCase(repo repository.ContatoRepo) *ContatoUseCase {
	return &ContatoUseCase{repo: repo}
}

func (u *ContatoUseCase) InitializeContact(contato models.Contact) (string, error) {
	count, err := u.repo.GetCountContatosAtivos()
	if err != nil {
		return "", apperror.Internal("Erro ao contar contatos ativos")
	}

	existente, err := u.repo.GetContatoByTelefone(contato.Telefone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if count >= config.GetPlanoLimite() {
				return "", apperror.Conflict("Limite do plano excedido")
			}
			if err := u.repo.CreateContato(contato.Nome, contato.Telefone); err != nil {
				return "", apperror.BadRequest(err.Error())
			}
			return "novo", nil
		}
		return "", apperror.Internal(err.Error())
	}

	conversationID := "NULL"
	if existente.ConversationID.Valid {
		conversationID = existente.ConversationID.String
	}
	return "veterano " + conversationID, nil
}

func (u *ContatoUseCase) SetConversationId(telefone, conversationId string) error {
	return u.repo.SetConversationId(telefone, conversationId)
}
