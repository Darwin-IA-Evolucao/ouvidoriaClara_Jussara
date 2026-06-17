package models

type Endereco struct {
	Logradouro string `json:"Logradouro"`
	Bairro     string `json:"Bairro"`
	Regiao     int    `json:"Região"`
}

type Logradouro struct {
	Logradouro string `db:"logradouro" json:"logradouro"`
	Bairro     string `db:"bairro" json:"bairro"`
	Regiao     string `db:"regiao" json:"regiao"`
}
