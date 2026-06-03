package models

type Usuario struct {
	Usuario string `json:"usuario" binding:"required" db:"usuario"`
	Senha   string `json:"senha" binding:"required" db:"senha"`
}
