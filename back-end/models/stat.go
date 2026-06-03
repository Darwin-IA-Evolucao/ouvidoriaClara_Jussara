package models

type Regiao struct {
	ID           int            `json:"id"`
	Distribuicao map[string]int `json:"distribuicao"`
}

type Stat struct {
	IndicacoesAprovadas    int      `json:"indicacoesAprovadas"`
	RequerimentosAprovados int      `json:"requerimentosAprovados"`
	OficiosAprovados       int      `json:"oficiosAprovados"`
	NumPessoas             int      `json:"numPessoas"`
	Indicacoes             int      `json:"indicacoes"`
	NumRequerimentos       int      `json:"numRequerimentos"`
	NumOficios             int      `json:"numOficios"`
	PercIndicacao          float64  `json:"percIndicacao"`
	NumReclamacoes         int      `json:"numReclamacoes"`
	NumIndicacoes          int      `json:"numIndicacoes"`
	NumReprovados          int      `json:"numReprovados"`
	Protocoladas           int      `json:"protocoladas"`
	PercProtocolacao       float64  `json:"percProtocolacao"`
	Regioes                []Regiao `json:"regioes"`
}
