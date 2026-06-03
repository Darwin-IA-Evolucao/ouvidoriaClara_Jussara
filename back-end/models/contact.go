package models

import (
	"database/sql"
	"time"
)

type Contact struct {
	ConversationID sql.NullString `db:"conversation_id" json:"conversation_id"`
	Nome           string         `db:"nome" json:"nome"`
	Telefone       string         `db:"telefone" json:"telefone"`
	Instance       *string        `db:"instance" json:"instance"`
	Ativo          bool           `db:"ativo" json:"ativo"`
	DataCriacao    string         `db:"data_criacao" json:"data_criacao"`
}

type Contato struct {
	Telefone       string    `json:"telefone" db:"telefone"`
	Nome           string    `json:"nome" db:"nome"`
	ConversationId *string   `json:"conversationId" db:"conversation_id"`
	Ativo          bool      `json:"ativo" db:"ativo"`
	Instance       *string   `json:"instance" db:"instance"`
	DataCriacao    time.Time `json:"data_criacao" db:"data_criacao"`
}

type GetConversationIdResponse struct {
	ConversationId string `json:"conversationId" db:"conversation_id"`
	Exists         bool   `json:"exists"`
}

type SetConversationIdRequest struct {
	Telefone       string `json:"telefone" db:"telefone"`
	ConversationId string `json:"conversationId" db:"conversation_id"`
}
