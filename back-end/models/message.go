package models

type Message struct {
	Telefone string `db:"telefone" json:"telefone"`
	Conteudo string `db:"conteudo" json:"conteudo"`
}

type GetMensagensResponse struct {
	TelefoneCliente string `json:"telefone_cliente"`
	Mensagem        string `json:"mensagem"`
	DataAtual       string `json:"data_atual"`
	HoraAtual       string `json:"hora_atual"`
	DiaSemana       string `json:"dia_semana"`
}
