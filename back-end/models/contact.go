package models

import (
	"database/sql"
	"time"
)

type Contact struct {
	ConversationID sql.NullString `db:"conversation_id" json:"conversation_id"`
	Nome           *string        `db:"nome" json:"nome"`
	Telefone       string         `db:"telefone" json:"telefone"`
	Instance       *string        `db:"instance" json:"instance"`
	Ativo          bool           `db:"ativo" json:"ativo"`
	DataCriacao    string         `db:"data_criacao" json:"data_criacao"`
	Campanha       *string        `json:"campanha" db:"campanha"`
}

type Contato struct {
	Telefone       string    `json:"telefone" db:"telefone"`
	Nome           *string   `json:"nome" db:"nome"`
	ConversationId *string   `json:"conversationId" db:"conversation_id"`
	Ativo          bool      `json:"ativo" db:"ativo"`
	Instance       *string   `json:"instance" db:"instance"`
	DataCriacao    time.Time `json:"data_criacao" db:"data_criacao"`
	Campanha       *string   `json:"campanha" db:"campanha"`
}

type GetConversationIdResponse struct {
	ConversationId string `json:"conversationId" db:"conversation_id"`
	Exists         bool   `json:"exists"`
}

type SetConversationIdRequest struct {
	Telefone       string `json:"telefone" db:"telefone"`
	ConversationId string `json:"conversationId" db:"conversation_id"`
}

type CreateContatoRequest struct {
	Telefone string `json:"telefone"`
	Nome     string `json:"nome"`
	Campanha string `json:"campanha"`
}

type ImportLinhaErro struct {
	Linha    int    `json:"linha"`
	Telefone string `json:"telefone"`
	Motivo   string `json:"motivo"`
}

type ImportContatosResult struct {
	Criados    int               `json:"criados"`
	Invalidos  []ImportLinhaErro `json:"invalidos"`
	Duplicados []ImportLinhaErro `json:"duplicados"`
}

type ContatosUnificados struct {
	Telefone       string  `json:"telefone" db:"telefone"`
	Nome           *string `json:"nome" db:"nome"`
	Cidade         *string `json:"cidade" db:"cidade"`
	Endereco       *string `json:"endereco" db:"endereco"`
	Bairro         *string `json:"bairro" db:"bairro"`
	DataNascimento *string `json:"dataNascimento" db:"data_nascimento"`
	DataCriacao    *string `json:"dataCriacao" db:"data_criacao"`
	ConversationId *string `json:"conversationId" db:"conversation_id"`
	Instance       *string `json:"instance" db:"instance"`
	Campanha       *string `json:"campanha" db:"campanha"`
	LeadAtivo      bool    `json:"leadAtivo" db:"lead_ativo"`
	DarwinAtivo    bool    `json:"darwinAtivo" db:"darwin_ativo"`
	IsGelado       bool    `json:"isGelado" db:"is_gelado"`
	IsCliente      bool    `json:"isCliente" db:"is_cliente"`
}

type ContatosUnificadosResponse struct {
	ContatosUnificados []ContatosUnificados `json:"contatos"`
	Total              int                  `json:"total"`
	Limite             int                  `json:"limite"`
	Usados             int                  `json:"usados"`
	Ocupacao           string               `json:"ocupacao"`
}
