package usecases

import (
	"back-end/models"
	"back-end/repository"
	"back-end/services"
	"log"
)

type ProtocoloUseCases struct {
	repository repository.ProtocoloRepository
}

func NewProtocoloUseCases(repo repository.ProtocoloRepository) ProtocoloUseCases {
	return ProtocoloUseCases{repository: repo}
}

func (uc ProtocoloUseCases) EnviarProtocolo(protocolo models.Protocolo) error {
	if err := uc.repository.CreateProtocolo(protocolo.Numero, protocolo.IDReclamacao); err != nil {
		return err
	}
	if err := uc.repository.MarcarResolvido(protocolo.IDReclamacao); err != nil {
		return err
	}

	dados, err := uc.repository.GetDadosReclamacao(protocolo.IDReclamacao)
	if err != nil {
		return err
	}

	if err := services.EnviarNotificacao(dados.Telefone, dados.Nome, dados.Categoria, protocolo.Numero); err != nil {
		log.Printf("Erro ao enviar notificação: %v", err)
	}

	if err := uc.repository.MarcarAvisado(protocolo.IDProtocolo); err != nil {
		log.Printf("Erro ao marcar protocolo como avisado: %v", err)
	}

	return nil
}
