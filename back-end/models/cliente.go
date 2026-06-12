package models

import (
	"fmt"
	"time"
)

type Cliente struct {
	Telefone       string    `db:"telefone" json:"telefone"`
	Nome           string    `db:"nome" json:"nome"`
	Cidade         string    `db:"cidade" json:"cidade"`
	Endereco       string    `db:"endereco" json:"endereco"`
	Bairro         string    `db:"bairro" json:"bairro"`
	DataNascimento string    `db:"data_nascimento" json:"dataNascimento"`
	DataCriacao    time.Time `db:"data_criacao" json:"dataCriacao"`
}

func (c Cliente) String() string {
	return fmt.Sprintf(
		`Cliente
Telefone: %s
Nome: %s
Cidade: %s
Endereço: %s
Bairro: %s
Data Nascimento: %s
Data Criação: %s`,
		c.Telefone,
		c.Nome,
		c.Cidade,
		c.Endereco,
		c.Bairro,
		c.DataNascimento,
		c.DataCriacao.Format("02/01/2006 15:04:05"),
	)
}
