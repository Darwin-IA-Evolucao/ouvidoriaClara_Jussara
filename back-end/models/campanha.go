package models

import (
	"back-end/utils"
	"strings"
)

type Campanha struct {
	ID           int    `json:"id" db:"id"`
	PalavraChave string `json:"palavraChave" db:"palavra_chave"`
	QtdContatos  int    `json:"qtdContatos" db:"qtd_contatos"`
}

type ContatoCampanha struct {
	Telefone string `json:"telefone" db:"telefone"`
	Nome     string `json:"nome" db:"nome"`
}

func (c *Campanha) NormalizePalavraChave() string {
	return utils.RemoveAcento(strings.TrimSpace(strings.ToLower(c.PalavraChave)))
}
