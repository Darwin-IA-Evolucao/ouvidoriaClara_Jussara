package models

type RequestData struct {
	Reclamacao string `json:"reclamacao" binding:"required"`
	Nome       string `json:"nome" binding:"required"`
	Telefone   string `json:"telefone" binding:"required"`
	Categoria  string `json:"categoria" binding:"required"`
	Regiao     int    `json:"regiao" binding:"required"`
}

type Inquerito struct {
	ID         int    `db:"id"`
	Reclamacao string `db:"reclamacao"`
	Nome       string `db:"nome"`
	Telefone   string `db:"telefone"`
	Categoria  string `db:"categoria"`
	Regiao     int    `db:"regiao"`
}

type Reclamacao struct {
	ID          int    `json:"id"`
	Nome        string `json:"nome"`
	Telefone    string `json:"telefone"`
	Categoria   string `json:"categoria"`
	Regiao      int    `json:"regiao"`
	Reclamacao  string `json:"reclamacao"`
	Resolvido   bool   `json:"resolvido"`
	DataCriacao string `json:"dataCriacao"`
	Tipo        string `json:"tipo"`
	Status      string `json:"status"`
	Protocolo   *int   `json:"num_protocolo"`
}

type Atualizacao struct {
	NovaReclamacao string `json:"novaReclamacao" binding:"required"`
}
