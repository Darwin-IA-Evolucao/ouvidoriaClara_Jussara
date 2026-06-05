package models

type Protocolo struct {
	IDProtocolo  int64 `json:"idProtocolo" db:"idprotocolo"`
	IDReclamacao int64 `json:"idReclamacao" db:"idreclamacao"`
	Numero       int64 `json:"numero" db:"numero"`
	Avisado      bool  `json:"avisado" db:"avisado"`
}
