package models

type Usuario struct {
	Celular string `db:"celular" json:"celular"`
	Senha   string `db:"senha" json:"senha"`
}
