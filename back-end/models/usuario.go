package models

type Usuario struct {
	ID      int    `db:"id" json:"id"`
	Celular string `db:"celular" json:"celular"`
	Senha   string `db:"senha" json:"senha"`
	Role    string `db:"role" json:"role"`
	Ativo   bool   `db:"ativo" json:"ativo"`
}
