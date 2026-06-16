package models

type Regiao struct {
	ID           int            `json:"id"`
	Distribuicao map[string]int `json:"distribuicao"`
}

type Stat struct {
	NumPessoas             int              `json:"numPessoas"`
	NumReclamacoes         int              `json:"numReclamacoes"`
	NumRequerimentos       int              `json:"numRequerimentos"`
	NumIndicacoes          int              `json:"numIndicacoes"`
	RequerimentosAprovados int              `json:"requerimentosAprovados"`
	IndicacoesAprovadas    int              `json:"indicacoesAprovadas"`
	NumReprovados          int              `json:"numReprovados"`
	Regioes                []StatsRegiao `json:"regioes"`
	Categorias             []StatsCategoria   `json:"categorias"`
	Tipos                  []StatsTipo      `json:"tipos"`
	PercIndicacao          float64          `json:"percIndicacao"`
}

type StatsCategoria struct {
	Categoria    string `json:"categoria" db:"categoria"`
	QtdCategoria int    `json:"qtdCategoria" db:"qtd_categoria"`
}

type StatsRegiao struct {
	Regiao    string `json:"regiao" db:"regiao"`
	QtdRegiao int    `json:"qtdRegiao" db:"qtd_regiao"`
}
type StatsTipo struct {
	Tipo      string `json:"tipo" db:"tipo"`
	QtdRegiao int    `json:"qtdTipo" db:"qtd_tipo"`
}
