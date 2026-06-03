package usecases

import (
	"back-end/apperror"
	"back-end/config"
	"back-end/models"
	"back-end/repository"
	"back-end/services"
	"sync"
	"time"
)

type MensagemUseCase struct {
	repo        repository.MensagemRepo
	respondendo bool
	mensagem    string
	counter     int
	mu          sync.Mutex
}

func NewMensagemUseCase(repo repository.MensagemRepo) *MensagemUseCase {
	return &MensagemUseCase{repo: repo}
}

func (u *MensagemUseCase) AddMensagem(msg models.Message) (string, error) {
	count, err := u.repo.GetCountContatosAtivos()
	if err != nil {
		return "", apperror.Internal("Erro ao contar contatos ativos")
	}

	existente, err := u.repo.GetCountContatoAtivoByTelefone(msg.Telefone)
	if err != nil {
		return "", apperror.Internal("Erro ao verificar contato existente")
	}

	if count >= config.GetPlanoLimite() && existente == 0 {
		return "", apperror.Conflict("Limite do plano excedido")
	}

	if err := u.repo.CreateMensagem(msg.Telefone, msg.Conteudo); err != nil {
		return "", apperror.BadRequest(err.Error())
	}

	data := time.Now().Format("02/01/2006")
	contador, err := u.repo.GetCountConversaHoje(msg.Telefone, data)
	if err != nil {
		return "", apperror.Internal(err.Error())
	}
	if contador < 1 {
		if err := u.repo.CreateConversa(msg.Telefone, data); err != nil {
			return "", apperror.BadRequest(err.Error())
		}
	}

	totalMensagens, err := u.repo.GetCountMensagens(msg.Telefone)
	if err != nil {
		return "", apperror.Internal(err.Error())
	}

	if err := u.repo.AtivarContato(msg.Telefone); err != nil {
		return "", apperror.BadRequest(err.Error())
	}

	u.mu.Lock()
	defer u.mu.Unlock()
	if totalMensagens > 1 {
		u.mensagem += msg.Conteudo + " "
		return "segunda", nil
	}
	u.respondendo = true
	return "primeira", nil
}

func (u *MensagemUseCase) GetMensagens(telefone string) (*models.GetMensagensResponse, error) {
	mensagens, err := u.repo.GetMensagensByTelefone(telefone)
	if err != nil {
		return nil, err
	}

	var finalMessage string
	for _, msg := range mensagens {
		finalMessage += msg.Conteudo + " "
	}

	now := time.Now()
	diaSemana := []string{"Domingo", "Segunda-feira", "Terça-feira", "Quarta-feira", "Quinta-feira", "Sexta-feira", "Sábado"}[now.Weekday()]

	return &models.GetMensagensResponse{
		TelefoneCliente: telefone,
		Mensagem:        finalMessage,
		DataAtual:       now.Format("02/01/2006"),
		HoraAtual:       now.Format("15:04"),
		DiaSemana:       diaSemana,
	}, nil
}

func (u *MensagemUseCase) ClearMensagens(telefone string) error {
	if err := u.repo.DeleteMensagens(telefone); err != nil {
		return apperror.BadRequest(err.Error())
	}

	conversationID, err := u.repo.GetConversationId(telefone)
	if err != nil {
		return apperror.Internal(err.Error())
	}

	u.mu.Lock()
	mensagem := u.mensagem
	respondendo := u.respondendo
	u.mu.Unlock()

	u.responder(telefone, conversationID, mensagem, respondendo)

	u.mu.Lock()
	u.respondendo = false
	u.mensagem = ""
	u.mu.Unlock()

	return nil
}

func (u *MensagemUseCase) responder(telefone, conversationID, mensagem string, respondendo bool) error {
	u.mu.Lock()
	u.counter++
	if u.counter > 1 {
		u.counter = 0
		u.mu.Unlock()
		return nil
	}
	u.mu.Unlock()

	if !respondendo || mensagem == "" {
		return nil
	}
	return services.Responder(telefone, conversationID, mensagem)
}
