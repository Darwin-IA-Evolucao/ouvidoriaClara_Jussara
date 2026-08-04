package models

import (
	"back-end/utils"
	"strings"
)

type Campanha struct {
	ID           int    `json:"id" db:"id"`
	PalavraChave string `json:"palavraChave" db:"palavra_chave"`
}

func (c *Campanha) NormalizePalavraChave() string {
	return utils.RemoveAcento(strings.TrimSpace(strings.ToLower(c.PalavraChave)))
}
