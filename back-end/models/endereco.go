package models

type Endereco struct {
	Logradouro string `json:"Logradouro"`
	Bairro     string `json:"Bairro"`
	Regiao     int    `json:"Região"`
}

type Logradouro struct {
	Logradouro string `db:"logradouro"`
	Bairro     string `db:"bairro"`
	Regiao     int    `db:"regiao"`
}
