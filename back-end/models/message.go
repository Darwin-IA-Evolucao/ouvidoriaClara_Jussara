package models

import (
	"sync"
	"time"
)

type Message struct {
	Id       int    `db:"id" json:"id"`
	Telefone string `db:"telefone" json:"telefone"`
	Conteudo string `db:"conteudo" json:"conteudo"`
}

type AtividadeCliente struct {
	Telefone        string `db:"telefone" json:"telefone"`
	UltimaInteracao string `db:"ultima_interacao" json:"ultimaInteracao"`
	Lembrete10Min   bool   `db:"lembrete_10min" json:"lembrete10Min"`
}

type Sessao struct {
	Respondendo bool
	Mu          sync.Mutex // Mutex para evitar processamento duplicado
}

// Mapa global thread-safe
var Sessoes sync.Map

type Mensagem struct {
	ID         int       `json:"id" db:"id"`
	Telefone   string    `json:"telefone" db:"telefone"`
	Conteudo   string    `json:"conteudo" db:"conteudo"`
	FoiEnviado bool      `json:"foiEnviado" db:"foienviado"`
	CreatedAt  time.Time `json:"createdAt" db:"data"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type AddMensagem struct {
	Telefone string `json:"telefone"`
	Nome     string `json:"nome"`
	Instance string `json:"instance"`
	Conteudo string `json:"conteudo"`
}

type GetMensagensResponse struct {
	Telefone  string `json:"telefone" db:"telefone"`
	Data      string `json:"data" db:"data"`
	Horario   string `json:"horario" db:"horario"`
	DiaSemana string `json:"diaSemana" db:"dia_semana"`
	Mensagem  string `json:"mensagem" db:"mensagem"`
}

type AvisoPlanoAtigido struct {
	Avisado         bool      `json:"avisado" db:"avisado"`
	TelefoneRicardo string    `json:"telefoneRicardo" db:"telefone_ricardo"`
	TelefoneLeo     string    `json:"telefoneLeo" db:"telefone_leo"`
	Data            time.Time `json:"data" db:"data"`
}