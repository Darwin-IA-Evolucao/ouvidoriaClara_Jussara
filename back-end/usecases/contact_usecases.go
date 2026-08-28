package usecases

import (
	"back-end/apperror"
	"back-end/models"
	"back-end/repository"
	"back-end/utils"
	"errors"
	"io"
	"strings"

	"github.com/lib/pq"
)

const (
	instanceOuvidoria = "ouvidoria_clara_jussara"
	campanhaSite      = "site"
	tamanhoMaxNome    = 120
)

type ContatoUseCase struct {
	repo         repository.ContatoRepo
	campanhaRepo *repository.CampanhaRepository
}

func NewContatoUseCase(repo repository.ContatoRepo, campanhaRepo *repository.CampanhaRepository) *ContatoUseCase {
	return &ContatoUseCase{repo: repo, campanhaRepo: campanhaRepo}
}

func (u *ContatoUseCase) GetAllContatos(limit, offset int) ([]models.Contact, error) {
	return u.repo.GetAllContatos(limit, offset)
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

func (u *ContatoUseCase) CreateContato(telefone, nome, campanha string) error {
	telefone = utils.NormalizeTelefone(telefone)
	if !utils.TelefoneValido(telefone) {
		return apperror.BadRequest("telefone em formato invalido")
	}
	campanha = strings.TrimSpace(campanha)
	if campanha == "" {
		return apperror.BadRequest("campanha obrigatoria")
	}
	instance := instanceOuvidoria
	nomeVal := strings.TrimSpace(nome)
	err := u.repo.CreateContato(&models.Contato{
		Telefone: telefone,
		Nome:     &nomeVal,
		Instance: &instance,
		Campanha: &campanha,
	})
	if err != nil {
		return mapContatoInsertErr(err)
	}
	return nil
}

// CreateLeadSite cadastra o lead vindo do site. Se o telefone ja existe, mantem o
// registro atual sem sobrescrever a campanha de origem.
func (u *ContatoUseCase) CreateLeadSite(telefone, nome, campanha string) error {
	tel := utils.NormalizeTelefone(telefone)
	if !utils.TelefoneValido(tel) {
		return apperror.BadRequest("telefone em formato invalido")
	}

	if campanha == "" {
		campanha = campanhaSite
	}
	
	if err := u.ensureCampanha(campanha); err != nil {
		return err
	}
	instance := instanceOuvidoria
	nomeVal := strings.TrimSpace(nome)
	if len(nomeVal) > tamanhoMaxNome {
		nomeVal = nomeVal[:tamanhoMaxNome]
	}
	return u.repo.CreateContatoIfNotExists(&models.Contato{
		Telefone: tel,
		Nome:     &nomeVal,
		Instance: &instance,
		Campanha: &campanha,
	})
}

func (u *ContatoUseCase) ImportContatos(filename string, file io.Reader, campanha string) (*models.ImportContatosResult, error) {
	c := &models.Campanha{PalavraChave: campanha}
	campanha = c.NormalizePalavraChave()
	if campanha == "" {
		return nil, apperror.BadRequest("campanha obrigatoria")
	}
	linhas, err := utils.ParsePlanilhaContatos(filename, file)
	if err != nil {
		return nil, apperror.BadRequest(err.Error())
	}
	result := &models.ImportContatosResult{
		Invalidos:  []models.ImportLinhaErro{},
		Duplicados: []models.ImportLinhaErro{},
	}
	instance := instanceOuvidoria
	for _, linha := range linhas {
		tel := utils.NormalizeTelefone(linha.Telefone)
		if !utils.TelefoneValido(tel) {
			result.Invalidos = append(result.Invalidos, models.ImportLinhaErro{
				Linha:    linha.Linha,
				Telefone: linha.Telefone,
				Motivo:   "telefone em formato invalido",
			})
			continue
		}
		nomeVal := strings.TrimSpace(linha.Nome)
		err := u.repo.CreateContato(&models.Contato{
			Telefone: tel,
			Nome:     &nomeVal,
			Instance: &instance,
			Campanha: &campanha,
		})
		if err != nil {
			if isUniqueViolation(err) {
				result.Duplicados = append(result.Duplicados, models.ImportLinhaErro{
					Linha:    linha.Linha,
					Telefone: tel,
					Motivo:   "telefone ja cadastrado",
				})
				continue
			}
			return nil, err
		}
		result.Criados++
	}
	if result.Criados > 0 {
		if err := u.ensureCampanha(campanha); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (u *ContatoUseCase) ensureCampanha(palavraChave string) error {
	c := &models.Campanha{PalavraChave: palavraChave}
	c.PalavraChave = c.NormalizePalavraChave()
	exists, err := u.campanhaRepo.ExistsCampanha(c.PalavraChave)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return u.campanhaRepo.CreateCampanha(c)
}

func mapContatoInsertErr(err error) error {
	if isUniqueViolation(err) {
		return apperror.Conflict("telefone ja cadastrado")
	}
	return err
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}
