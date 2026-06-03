package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	_ "github.com/lib/pq" // Import necesario para o driver do PostgreSQL
	"gopkg.in/gomail.v2"
)

// Estrutura para receber os dados JSON da requisição
type RequestData struct {
	Reclamacao string `json:"reclamacao" binding:"required"`
	Nome       string `json:"nome" binding:"required"`
	Telefone   string `json:"telefone" binding:"required"`
	Categoria  string `json:"categoria" binding:"required"`
	Regiao     int    `json:"regiao" binding:"required"`
}

type Inquerito struct {
	ID         int
	Reclamacao string
	Nome       string
	Telefone   string
	Categoria  string
	Regiao     int
}

type Contact struct {
	ConversationID sql.NullString `json:"conversation_id"`
	Nome           string         `json:"nome"`
	Telefone       string         `json:"telefone"`
}

type Message struct {
	Telefone string `json:"telefone"`
	Conteudo string `json:"conteudo"`
}

type contactSetMold struct {
	ConversationID string `json:"conversation_id"`
	Telefone       string `json:"telefone"`
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

type Logradouro struct {
	Logradouro string
	Bairro     string
	Regiao     int
}

type Protocolo struct {
	ID     int64 `json:"id"`     // idReclamacao
	Numero int64 `json:"numero"` // número do protocolo
}

type Usuario struct {
	Usuario string `json:"usuario" binding:"required"`
	Senha   string `json:"senha" binding:"required"`
}

type Atualizacao struct {
	NovaReclamacao string `json:"novaReclamacao" binding:"required"`
}

type Endereco struct {
	Logradouro string `json:"Logradouro"`
	Bairro     string `json:"Bairro"`
	Regiao     int    `json:"Região"`
}

type Mensagem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Index        int      `json:"index"`
	Message      Mensagem `json:"message"` // Alterado para string
	FinishReason string   `json:"finish_reason"`
}

type Lead struct {
	Nome     string
	Telefone string
	Ativo    bool
}

var db *sql.DB
var respondendo bool
var mensagem string
var counter int
var categorias = []string{"saúde", "educação", "transporte", "asfalto", "governança"}
var planos = map[string]int{"bronze": 200, "prata": 1000, "ouro": 5000}

// Constantes para horários de funcionamento
const (
	plano_atual = "ouro"
)

func normalizeTelefone(telefone string) string {
	// Remove todos os caracteres que não são números
	re := regexp.MustCompile(`\D`)
	somenteNumeros := re.ReplaceAllString(telefone, "")

	// Remove o zero inicial, se houver
	if strings.HasPrefix(somenteNumeros, "0") {
		somenteNumeros = somenteNumeros[1:]
	}

	// Adiciona o prefixo 55 se ainda não estiver presente
	if !strings.HasPrefix(somenteNumeros, "55") {
		somenteNumeros = "55" + somenteNumeros
	}

	// Garante que o número tenha no máximo 13 caracteres (55 + DDD + número)
	if len(somenteNumeros) > 13 {
		somenteNumeros = somenteNumeros[:13]
	}

	return somenteNumeros
}

func initializeContact(c *gin.Context) {
	var contact Contact
	var verCont Contact

	if err := c.BindJSON(&contact); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	if contact.Telefone == "status" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Telefone obrigatório."})
		return
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM contatos WHERE ativo = true").Scan(&count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar contatos ativos"})
		return
	}

	row := db.QueryRow("SELECT nome, telefone, conversation_id FROM contatos WHERE telefone = $1", contact.Telefone)
	if err := row.Scan(&verCont.Nome, &verCont.Telefone, &verCont.ConversationID); err != nil {
		if err == sql.ErrNoRows {
			if count >= planos[plano_atual] {
				c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Limite do plano excedido"})
				return
			}
			_, err := db.Exec("INSERT INTO contatos (nome, telefone, conversation_id, ativo) VALUES ($1, $2, $3, true)", contact.Nome, contact.Telefone, nil)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.String(http.StatusCreated, "novo")
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	conversationID := "NULL"
	if verCont.ConversationID.Valid {
		conversationID = verCont.ConversationID.String
	}
	c.String(http.StatusOK, "veterano "+conversationID)
}

func addMessage(c *gin.Context) {
	var msg Message

	if err := c.BindJSON(&msg); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	var count, existingContact int
	err := db.QueryRow("SELECT COUNT(*) FROM contatos WHERE ativo = true").Scan(&count)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar contatos ativos"})
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM contatos WHERE telefone = $1 AND ativo = true", msg.Telefone).Scan(&existingContact)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao verificar contato existente"})
		return
	}

	if count >= planos[plano_atual] && existingContact == 0 {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "Limite do plano excedido"})
		return
	}
	_, err = db.Exec("INSERT INTO mensagens (telefone, conteudo) VALUES ($1, $2)", msg.Telefone, msg.Conteudo)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hoje := time.Now()
	data := hoje.Format("02/01/2006")
	var contador int64
	iters, dbErr := db.Query("SELECT * FROM Conversas WHERE telefone = $1 AND data = $2", msg.Telefone, data)
	if dbErr != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
		return
	}
	defer iters.Close()
	for iters.Next() {
		contador++
	}
	if contador < 1 {
		_, outroErr := db.Exec("INSERT INTO Conversas (telefone, data) VALUES ($1, $2)", msg.Telefone, data)
		if outroErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": outroErr.Error()})
			return
		}
	}

	rows, err := db.Query("SELECT * FROM mensagens WHERE telefone = $1", msg.Telefone)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, dbErr = db.Exec("UPDATE Contatos SET ativo = true WHERE telefone = $1", msg.Telefone)
	if dbErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	defer rows.Close()

	var counter int64
	for rows.Next() {
		counter++
	}

	if counter > 1 {
		mensagem += msg.Conteudo + " "
		c.String(http.StatusOK, "segunda")
		return
	}
	respondendo = true
	c.String(http.StatusOK, "primeira")
}

func setConversation(c *gin.Context) {

	var mold contactSetMold
	if err := c.BindJSON(&mold); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	_, dbErr := db.Exec("UPDATE Contatos SET conversation_id = $1 WHERE telefone = $2", mold.ConversationID, mold.Telefone)
	if dbErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	c.String(http.StatusOK, "Sucesso!")
}

func getMessagesByPhone(c *gin.Context) {
	tel := c.Param("telefone")
	var finalMessage string
	now := time.Now()
	dataAtual := now.Format("02/01/2006")
	hora := now.Format("15:04")
	diaSemana := []string{"Domingo", "Segunda-feira", "Terça-feira", "Quarta-feira", "Quinta-feira", "Sexta-feira", "Sábado"}[now.Weekday()]

	rows, err := db.Query("SELECT telefone, conteudo FROM mensagens WHERE telefone = $1", tel)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msg Message
		if dbErr := rows.Scan(&msg.Telefone, &msg.Conteudo); dbErr != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": dbErr.Error()})
			return
		}
		finalMessage += msg.Conteudo + " "
	}

	c.IndentedJSON(http.StatusOK, gin.H{"telefone_cliente": tel, "mensagem": finalMessage, "data_atual": dataAtual, "hora_atual": hora, "dia_semana": diaSemana})
}

func clearMessagesByPhone(c *gin.Context) {

	tel := c.Param("telefone")
	_, dbErr := db.Exec("DELETE FROM mensagens WHERE telefone =  $1", tel)
	if dbErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": dbErr.Error()})
		return
	}
	var conversationID string
	row := db.QueryRow("SELECT conversation_id FROM contatos WHERE telefone = $1", tel)
	if err := row.Scan(&conversationID); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Sucesso!")
	responder(tel, conversationID, mensagem)
	respondendo = false
	mensagem = ""
}

func responder(telefone, conversation_id, mensagem string) error {
	counter++
	if counter > 1 {
		counter = 0
		return nil
	}
	if !respondendo || mensagem == "" {
		return nil
	}
	baseURL := "http://147.93.10.167:5678/webhook/92232d3b-f9c8-4525-827e-0aa9fe181e27Arruda"
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Use a função auxiliar
	telefonePadronizado := padronizaTelefone(telefone)

	data := url.Values{}
	data.Set("user", telefonePadronizado)
	data.Set("data", mensagem)
	data.Set("telefone", telefonePadronizado)
	data.Set("conversation_id", conversation_id)

	req, err := http.NewRequest("POST", baseURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao enviar requisição: %s", resp.Status)
	}
	return nil
}

// Mensagem de boas-vindas
func greetings(c *gin.Context) {
	c.String(http.StatusOK, "Jesus é o Senhor!")
}

func generatePDF(requestData Inquerito) error {
	// Inicializa um novo PDF A4 em modo retrato
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Configurar UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Adicionar página
	pdf.AddPage()

	// Margens
	pdf.SetMargins(20, 20, 20)

	// Criar pasta para imagens se não existir
	if _, err := os.Stat("img"); os.IsNotExist(err) {
		os.Mkdir("img", 0755)
	}

	// Verifica se o logo existe
	logoFile := "/root/arruda_Clara/img/logo.png"
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		log.Printf("Atenção: Logo não encontrado em %s", logoFile)
		// Se quiser você pode criar um logo placeholder aqui
	} else {
		// Adicionar logo
		pdf.ImageOptions(
			logoFile,
			20,    // x position
			15,    // y position
			170,   // width
			30,    // height
			false, // flow
			gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
			0,
			"",
		)
	}

	// Cabeçalho - Título centralizado em negrito
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(60, 50)
	pdf.Cell(90, 10, tr("GABINETE DO VEREADOR MARCOS ARRUDA"))

	// Linha fina abaixo do título
	pdf.SetXY(85, 60)
	// Subtítulo "INDICAÇÃO N.º____ /2025"
	pdf.SetFont("Arial", "B", 12)
	indicacaoTitulo := "INDICAÇÃO N.º " + "___" + " / 2025"
	pdf.Cell(0, 10, tr(indicacaoTitulo))

	// Espaço após cabeçalho
	pdf.Ln(30)

	// Saudação "Senhor Presidente,"
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Excelentíssimo Senhor Presidente,"))
	pdf.Ln(15)

	// Texto da indicação
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" INDICA"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" ao Senhor Prefeito Municipal, que "))

	// Texto da reclamação (em texto normal)
	pdf.SetFont("Arial", "", 11)
	// Usando MultiCell com largura definida para 0 (que ocupa toda a largura disponível)
	pdf.Write(10, tr(requestData.Reclamacao+"."))
	pdf.Ln(15) // Espaço após a reclamação

	// Fechamento
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Nestes termos,"))
	pdf.Ln(5)
	pdf.Cell(0, 10, tr("Aguarda deferimento."))
	pdf.Ln(30)

	now := time.Now()
	// Formata a data por extenso (em português) - ex: 13 de março de 2025
	meses := []string{
		"janeiro", "fevereiro", "março", "abril", "maio", "junho",
		"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}
	dataStr := fmt.Sprintf("São Roque, %d de %s de %d", now.Day(), meses[now.Month()-1], now.Year())
	// Data
	pdf.Cell(0, 10, tr(dataStr))
	pdf.Ln(20)

	// Assinatura
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	// Salvar o arquivo PDF
	return pdf.OutputFileAndClose("indicacao.pdf")
}

func generateRequerimento(requestData Inquerito) error {
	// Garante que a string começa com "$$"
	if !strings.Contains(requestData.Reclamacao, "$$") {
		requestData.Reclamacao = "$$" + requestData.Reclamacao
	}
	args := strings.Split(requestData.Reclamacao, "$$")
	// Inicializa um novo PDF A4 em modo retrato
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Configurar UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Adicionar página
	pdf.AddPage()

	// Margens
	pdf.SetMargins(20, 20, 20)

	// Criar pasta para imagens se não existir
	if _, err := os.Stat("img"); os.IsNotExist(err) {
		os.Mkdir("img", 0755)
	}

	// Verifica se o logo existe
	logoFile := "/root/arruda_Clara/img/logo.png"
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		log.Printf("Atenção: Logo não encontrado em %s", logoFile)
		// Se quiser você pode criar um logo placeholder aqui
	} else {
		// Adicionar logo
		pdf.ImageOptions(
			logoFile,
			20,    // x position
			15,    // y position
			170,   // width
			30,    // height
			false, // flow
			gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
			0,
			"",
		)
	}

	// Cabeçalho - Título centralizado em negrito
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(60, 50)
	pdf.Cell(90, 10, tr("GABINETE DO VEREADOR MARCOS ARRUDA"))

	// Linha fina abaixo do título
	pdf.SetXY(85, 60)
	// Subtítulo "INDICAÇÃO N.º____ /2025"
	pdf.SetFont("Arial", "B", 12)
	requerimentoTitulo := "REQUERIMENTO N.º " + "___" + " / 2025"
	pdf.Cell(0, 10, tr(requerimentoTitulo))

	// Espaço após cabeçalho
	pdf.Ln(30)

	// Saudação "Senhor Presidente,"
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Senhor Presidente,"))
	pdf.Ln(15)

	// Saudação "Senhores Vereadores,"
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Senhores Vereadores,"))
	pdf.Ln(15)

	// Texto da indicação
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" REQUER"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" que, seja oficiado ao Prefeito Municipal para que nos sejam prestadas as seguintes informações:"))
	pdf.Ln(15)

	requests := strings.Split(args[1], "|")
	var count int
	for i, request := range requests {
		count = i + 1
		// Texto da reclamação (em texto normal)
		pdf.SetFont("Arial", "", 11)
		// Usando MultiCell com largura definida para 0 (que ocupa toda a largura disponível)
		pdf.Write(10, tr(fmt.Sprintf("%d. %s", count, request)))
		pdf.Ln(15) // Espaço após a reclamação
	}
	now := time.Now()
	// Formata a data por extenso (em português) - ex: 13 de março de 2025
	meses := []string{
		"janeiro", "fevereiro", "março", "abril", "maio", "junho",
		"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}
	dataStr := fmt.Sprintf("São Roque, %d de %s de %d", now.Day(), meses[now.Month()-1], now.Year())
	// Data
	pdf.Cell(0, 10, tr(dataStr))
	pdf.Ln(20)

	// Assinatura
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	// Salvar o arquivo PDF
	return pdf.OutputFileAndClose("requerimento.pdf")
}

func ehCategoria(str string) bool {
	str = strings.ToLower(str)
	for _, cat := range categorias {
		if str == cat {
			return true
		}
	}
	return false
}

func generateOficio(requestData Inquerito) error {
	// Inicializa um novo PDF A4 em modo retrato
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Configurar UTF-8
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	// Adicionar página
	pdf.AddPage()

	// Margens
	pdf.SetMargins(20, 20, 20)

	// Criar pasta para imagens se não existir
	if _, err := os.Stat("img"); os.IsNotExist(err) {
		os.Mkdir("img", 0755)
	}

	// Verifica se o logo existe
	logoFile := "/root/arruda_Clara/img/logo.png"
	if _, err := os.Stat(logoFile); os.IsNotExist(err) {
		log.Printf("Atenção: Logo não encontrado em %s", logoFile)
		// Se quiser você pode criar um logo placeholder aqui
	} else {
		// Adicionar logo
		pdf.ImageOptions(
			logoFile,
			20,    // x position
			15,    // y position
			170,   // width
			30,    // height
			false, // flow
			gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
			0,
			"",
		)
	}
	now := time.Now()
	// Cabeçalho - Título centralizado em negrito
	pdf.SetFont("Arial", "B", 12)
	pdf.SetXY(60, 50)
	pdf.Cell(90, 10, tr("GABINETE DO VEREADOR MARCOS ARRUDA"))

	// Linha fina abaixo do título
	pdf.SetXY(85, 60)
	// Subtítulo "INDICAÇÃO N.º____ /2025"
	pdf.SetFont("Arial", "B", 12)
	indicacaoTitulo := "OFÍCIO N.º " + "___" + " / " + now.Format("2006")
	pdf.Cell(0, 10, tr(indicacaoTitulo))

	// Espaço após cabeçalho
	pdf.Ln(30)

	// Saudação "Senhor Presidente,"
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Senhor Presidente,"))
	pdf.Ln(15)

	// Texto da indicação
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr("Este Vereador, submentendo este documento ao Chefe do Poder Executivo, diretamente ou através de departamento ou divisão competente,"))
	pdf.SetFont("Arial", "B", 11)
	pdf.Write(10, tr(" OFICIA"))
	pdf.SetFont("Arial", "", 11)
	pdf.Write(10, tr(" ao Senhor Prefeito Municipal, que "))

	// Texto da reclamação (em texto normal)
	pdf.SetFont("Arial", "", 11)
	// Usando MultiCell com largura definida para 0 (que ocupa toda a largura disponível)
	pdf.Write(10, tr(requestData.Reclamacao+"."))
	pdf.Ln(15) // Espaço após a reclamação

	// Fechamento
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 10, tr("Nestes termos,"))
	pdf.Ln(5)
	pdf.Cell(0, 10, tr("Aguarda deferimento."))
	pdf.Ln(30)

	// Formata a data por extenso (em português) - ex: 13 de março de 2025
	meses := []string{
		"janeiro", "fevereiro", "março", "abril", "maio", "junho",
		"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
	}
	dataStr := fmt.Sprintf("São Roque, %d de %s de %d", now.Day(), meses[now.Month()-1], now.Year())
	// Data
	pdf.Cell(0, 10, tr(dataStr))
	pdf.Ln(20)

	// Assinatura
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 10, tr("VEREADOR MARCOS ARRUDA"), "", 1, "C", false, 0, "")

	// Salvar o arquivo PDF
	return pdf.OutputFileAndClose("oficio.pdf")
}

// Função para enviar e-mail com o PDF anexado
func sendEmail(to, subject, body, filename string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "gabinetevereadorarruda@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(filename) // Anexa o PDF gerado
	d := gomail.NewDialer("smtp.gmail.com", 587, "gabinetevereadorarruda@gmail.com", "jhbn cjnn eait roja")
	return d.DialAndSend(m)
}

func messageHandler(c *gin.Context) {
	var data RequestData
	// Bind JSON e valida os campos obrigatórios
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos são obrigatórios"})
		return
	}

	if !ehCategoria(data.Categoria) {
		c.JSON(http.StatusOK, gin.H{"error": "Categoria inválida"})
		return
	}

	// Normaliza o telefone
	data.Telefone = normalizeTelefone(data.Telefone)
	unwantedPhrases := []string{
		"Indico ao Excelentíssimo Senhor Prefeito Municipal que,",
		"Considerando que",
	}

	for _, phrase := range unwantedPhrases {
		data.Reclamacao = strings.TrimSpace(strings.ReplaceAll(data.Reclamacao, phrase, ""))
	}

	data.Categoria = strings.ToLower(data.Categoria)
	_, err := db.Exec("INSERT INTO Reclamacao (nome, telefone, categoria, regiao, reclamacao) VALUES ($1, $2, $3, $4, $5)", data.Nome, data.Telefone, data.Categoria, data.Regiao, data.Reclamacao)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Solicitação enviada com sucesso!"})
}

func enviaInquerito(data Inquerito) (bool, error) {
	filename := "indicacao.pdf"
	err := generatePDF(data)
	if err != nil {
		log.Printf("Erro ao gerar PDF: %v", err)
		return false, err
	}

	// Verificar se o arquivo existe após a geração
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Arquivo PDF não foi criado corretamente: %v", err)
		return false, err
	}

	err = sendEmail("gabinetevereadorarruda@gmail.com", "indicação", "Segue em anexo uma indicação requisitada por um munícipe de São Roque, obrigado!.", filename)
	if err != nil {
		log.Printf("Erro ao enviar email: %v", err)
		return false, err
	}

	log.Printf("Email enviado com sucesso para gabinetevereadorarruda@gmail.com")

	/*err = sendEmail("gabrielpadawan912@gmail.com", "inquerito", "Segue em anexo um inquerito requisitado por um munícipe de São Roque, obrigado!.", filename)
	  if err != nil {
	      log.Printf("Erro ao enviar email: %v", err)
	      return false, err
	  }

	  log.Printf("Email enviado com sucesso para gabrielpadawan912@gmail.com")
	*/
	// Apaga o arquivo PDF após o envio para liberar espaço na VPS
	os.Remove(filename)

	return true, nil
}

func enviaRequerimento(data Inquerito) (bool, error) {
	filename := "requerimento.pdf"
	err := generateRequerimento(data)
	if err != nil {
		log.Printf("Erro ao gerar PDF: %v", err)
		return false, err
	}

	// Verificar se o arquivo existe após a geração
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Arquivo PDF não foi criado corretamente: %v", err)
		return false, err
	}

	err = sendEmail("gabinetevereadorarruda@gmail.com", "requerimento", "Segue em anexo um requerimento requisitado por um munícipe de São Roque, obrigado!.", filename)
	if err != nil {
		log.Printf("Erro ao enviar email: %v", err)
		return false, err
	}

	log.Printf("Email enviado com sucesso para gabinetevereadorarruda@gmail.com")

	// Apaga o arquivo PDF após o envio para liberar espaço na VPS
	os.Remove(filename)

	return true, nil
}

func enviaOficio(data Inquerito) (bool, error) {
	filename := "oficio.pdf"
	err := generateOficio(data)
	if err != nil {
		log.Printf("Erro ao gerar PDF: %v", err)
		return false, err
	}

	// Verificar se o arquivo existe após a geração
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Arquivo PDF não foi criado corretamente: %v", err)
		return false, err
	}

	err = sendEmail("gabinetevereadorarruda@gmail.com", "ofício", "Segue em anexo um ofício requisitado por um munícipe de São Roque, obrigado!.", filename)
	if err != nil {
		log.Printf("Erro ao enviar email: %v", err)
		return false, err
	}

	log.Printf("Email enviado com sucesso para gabinetevereadorarruda@gmail.com")

	// Apaga o arquivo PDF após o envio para liberar espaço na VPS
	os.Remove(filename)

	return true, nil
}

func aprovarComoAmbos(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE Reclamacao SET status = 'aprovado' WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var data Inquerito
	row := db.QueryRow("SELECT id, reclamacao, nome, telefone, categoria, regiao FROM reclamacao WHERE id = $1", id)
	if err := row.Scan(&data.ID, &data.Reclamacao, &data.Nome, &data.Telefone, &data.Categoria, &data.Regiao); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !strings.Contains(data.Reclamacao, "$$") {
		if _, err := enviaInquerito(data); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, convErr := convertIndicacao(data.Reclamacao)
		if convErr != nil {
			c.String(http.StatusBadRequest, "Deu ruim "+convErr.Error())
			return
		}
		data.Reclamacao = body
		if _, err := enviaRequerimento(data); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if strings.Contains(data.Reclamacao, "$$") {
		if _, err := enviaRequerimento(data); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, convErr := convertRequerimento(data.Reclamacao)
		if convErr != nil {
			c.String(http.StatusBadRequest, "Deu ruim "+convErr.Error())
			return
		}
		data.Reclamacao = body
		if _, err := enviaInquerito(data); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inqueritos aprovados e enviados com sucesso!"})
}

func aprovarInquerito(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE Reclamacao SET status = 'aprovado' WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var data Inquerito
	row := db.QueryRow("SELECT id, reclamacao, nome, telefone, categoria, regiao FROM reclamacao WHERE id = $1", id)
	if err := row.Scan(&data.ID, &data.Reclamacao, &data.Nome, &data.Telefone, &data.Categoria, &data.Regiao); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := enviaInquerito(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito aprovado e enviado com sucesso!"})
}

func aprovarRequerimento(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE Reclamacao SET status = 'aprovado', tipo = 'requerimento' WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var data Inquerito
	row := db.QueryRow("SELECT id, reclamacao, nome, telefone, categoria, regiao FROM reclamacao WHERE id = $1", id)
	if err := row.Scan(&data.ID, &data.Reclamacao, &data.Nome, &data.Telefone, &data.Categoria, &data.Regiao); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := enviaRequerimento(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Requerimento aprovado e enviado com sucesso!"})
}

func aprovarOficio(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE Reclamacao SET status = 'aprovado', tipo = 'ofício' WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var data Inquerito
	row := db.QueryRow("SELECT id, reclamacao, nome, telefone, categoria, regiao FROM reclamacao WHERE id = $1", id)
	if err := row.Scan(&data.ID, &data.Reclamacao, &data.Nome, &data.Telefone, &data.Categoria, &data.Regiao); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := enviaOficio(data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ofício aprovado e enviado com sucesso!"})
}

func reprovarInquerito(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE Reclamacao SET status = 'reprovado' WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Inquerito reprovado com sucesso!"})
}

func editarReclamacao(c *gin.Context) {
	id := c.Param("id")
	var data Atualizacao
	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE reclamacao SET reclamacao = $1 WHERE id = $2", data.NovaReclamacao, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reclamação editada com sucesso!"})
}

func getAllReclamacoes(c *gin.Context) {
	rows, err := db.Query(`
		SELECT r.id, r.nome, r.telefone, r.categoria, r.regiao, r.resolvido, r.data, r.reclamacao, r.status, r.tipo, p.numero AS protocolo
		FROM reclamacao r
		LEFT JOIN protocolo p ON r.id = p.idReclamacao
		ORDER BY data DESC;
	`)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var reclamacoes []Reclamacao

	for rows.Next() {
		var reclamacao Reclamacao
		var protocolo sql.NullInt64 // Usando sql.NullInt64 para lidar com valores nulos

		if err := rows.Scan(&reclamacao.ID, &reclamacao.Nome, &reclamacao.Telefone, &reclamacao.Categoria,
			&reclamacao.Regiao, &reclamacao.Resolvido, &reclamacao.DataCriacao, &reclamacao.Reclamacao, &reclamacao.Status, &reclamacao.Tipo, &protocolo); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Verificando se o protocolo é NULL e atribuindo o valor (ou nil) ao campo
		if protocolo.Valid {
			reclamacao.Protocolo = new(int)              // Cria um novo ponteiro para int
			*reclamacao.Protocolo = int(protocolo.Int64) // Atribui o valor ao ponteiro
		} else {
			reclamacao.Protocolo = nil // Define como nil se for NULL no banco de dados
		}
		reclamacao.Categoria = strings.ToLower(reclamacao.Categoria)
		reclamacoes = append(reclamacoes, reclamacao)
	}

	c.IndentedJSON(http.StatusOK, reclamacoes)
}

func buscaBinaria(lista []Logradouro, ruaProc string) (int, bool) {
	var ruaAtual string
	ruaProc = strings.ToLower(ruaProc)
	inicio := 0
	fim := len(lista) - 1
	meio := (fim + inicio) / 2

	for fim >= inicio {
		ruaAtual = strings.ToLower(lista[meio].Logradouro)
		if ruaAtual == ruaProc {
			return meio, true // Rua encontrada
		}
		if ruaAtual > ruaProc {
			fim = meio - 1
			meio = inicio + (fim-inicio)/2
		} else {
			inicio = meio + 1
			meio = inicio + (fim-inicio)/2
		}
	}
	return -1, false // Rua não encontrada
}

// LevenshteinDistance calcula a distância de edição entre duas strings
// Esta distância representa o número mínimo de operações (inserções,
// exclusões ou substituições) necessárias para transformar uma string em outra
func LevenshteinDistance(s1, s2 string) int {
	// Converte para minúsculas para comparação sem distinção de maiúsculas/minúsculas
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	// Cria a matriz para calcular a distância
	rows := len(s1) + 1
	cols := len(s2) + 1
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		matrix[i][0] = i
	}

	for j := 1; j < cols; j++ {
		matrix[0][j] = j
	}

	// Preenche a matriz calculando a distância
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			// Escolhe a operação de menor custo (inserção, exclusão ou substituição)
			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deleção
				matrix[i][j-1]+1,      // inserção
				matrix[i-1][j-1]+cost, // substituição
			)
		}
	}

	// Retorna o valor final que representa a distância de edição
	return matrix[rows-1][cols-1]
}

// min3 retorna o menor de três inteiros
func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// EncontraMelhorCorrespondencia busca a string mais próxima na lista com base na distância de Levenshtein
// e também verifica se o termo de busca é uma substring do nome da rua
func EncontraMelhorCorrespondencia(lista []Logradouro, termo string, limiteSimilaridade int) (int, bool) {
	termoNormalizado := strings.ToLower(strings.TrimSpace(termo))

	// Se o termo está vazio, não pode corresponder a nada
	if termoNormalizado == "" {
		return -1, false
	}

	// Aumenta o limite de similaridade para permitir correspondências mais flexíveis
	if limiteSimilaridade <= 0 {
		// Define um limite mais generoso
		limiteSimilaridade = len(termo)
		if limiteSimilaridade > 10 {
			limiteSimilaridade = 10 // Máximo de 10 caracteres de diferença
		}
	}

	melhorIndice := -1
	melhorPontuacao := -1 // Quanto menor a pontuação, melhor a correspondência

	for i, item := range lista {
		logradouroNormalizado := strings.ToLower(strings.TrimSpace(item.Logradouro))

		// Inicializa com a distância de Levenshtein
		distancia := LevenshteinDistance(termoNormalizado, logradouroNormalizado)

		// Calcula uma pontuação composta (quanto menor, melhor)
		pontuacao := distancia

		// Verifica se o termo de busca é uma substring do logradouro
		// Isso dá prioridade a casos onde o termo está contido no nome da rua
		if strings.Contains(logradouroNormalizado, termoNormalizado) {
			// Dá uma grande vantagem para correspondências de substring
			pontuacao = pontuacao / 2
		}

		// Verifica correspondência de palavras-chave individuais
		// Por exemplo, se busca "Júlio Prestes", verifica se as palavras "júlio" e "prestes"
		// aparecem no nome do logradouro
		termoPalavras := strings.Fields(termoNormalizado)
		palavrasEncontradas := 0

		for _, palavra := range termoPalavras {
			if strings.Contains(logradouroNormalizado, palavra) {
				palavrasEncontradas++
			}
		}

		// Se todas as palavras-chave estiverem presentes, melhora muito a pontuação
		if palavrasEncontradas == len(termoPalavras) && len(termoPalavras) > 0 {
			pontuacao = pontuacao / 3
		} else if palavrasEncontradas > 0 {
			// Melhora parcialmente a pontuação com base na proporção de palavras encontradas
			pontuacao = pontuacao * (len(termoPalavras) - palavrasEncontradas + 1) / len(termoPalavras)
		}

		// Atualiza o melhor resultado se encontrar uma pontuação melhor
		if melhorPontuacao == -1 || pontuacao < melhorPontuacao {
			melhorPontuacao = pontuacao
			melhorIndice = i
		}
	}

	// Uma pontuação muito baixa é um bom sinal de correspondência
	// mesmo que a distância original seja alta
	if melhorPontuacao <= limiteSimilaridade || melhorPontuacao < 5 {
		return melhorIndice, true
	}

	return -1, false
}

func limparEndereco(input string) string {
	frasesIrrelevantes := []string{
		"em frente ao número", "altura do número", "altura do", "próximo à",
		"próximo ao", "final do", "início da", "em frente", "perto de",
		"perto do", "ao lado de", "na esquina com", "número", "prox",
	}
	input = strings.ToLower(input)
	for _, frase := range frasesIrrelevantes {
		if idx := strings.Index(input, frase); idx != -1 {
			return strings.TrimSpace(input[:idx])
		}
	}
	// Também limpa sufixos após vírgula
	input = strings.Split(input, ",")[0]
	return strings.TrimSpace(input)
}

func getRegiao(c *gin.Context) {
	input := c.Query("rua")

	// Etapa 1: limpar partes irrelevantes do endereço
	input = limparEndereco(input)

	// Etapa 2: tentativa rápida com ILIKE no banco
	var regiaoSQL int
	sqlErr := db.QueryRow(`
		SELECT regiao FROM enderecos 
		WHERE logradouro ILIKE '%' || $1 || '%' 
		ORDER BY LENGTH(logradouro) ASC LIMIT 1;
	`, input).Scan(&regiaoSQL)

	if sqlErr == nil {
		c.JSON(http.StatusOK, regiaoSQL)
		return
	}

	// Etapa 3: Carrega todos os logradouros (como na versão original)
	rows, err := db.Query("SELECT * FROM enderecos ORDER BY logradouro;")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	logradouros := []Logradouro{}
	for rows.Next() {
		var lg Logradouro
		if err := rows.Scan(&lg.Logradouro, &lg.Bairro, &lg.Regiao); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logradouros = append(logradouros, lg)
	}

	// Etapa 4: Busca binária exata
	index, achou := buscaBinaria(logradouros, input)
	if achou {
		c.JSON(200, logradouros[index].Regiao)
		return
	}

	// Etapa 5: Busca aproximada com Levenshtein
	limiteSimilaridade := len(input)
	if limiteSimilaridade > 10 {
		limiteSimilaridade = 10
	}

	melhorIndex, encontrado := EncontraMelhorCorrespondencia(logradouros, input, limiteSimilaridade)
	if encontrado {
		c.JSON(200, logradouros[melhorIndex].Regiao)
		return
	}

	// Etapa final: fallback
	c.JSON(200, 1)
}

func getStats(c *gin.Context) {
	// First query
	rows, err := db.Query("SELECT " +
		"(SELECT COUNT(DISTINCT telefone) FROM reclamacao) AS pessoas_com_reclamacao, " +
		"(SELECT COUNT(DISTINCT telefone) FROM contatos " +
		"WHERE telefone NOT IN (SELECT telefone FROM reclamacao)" +
		") AS pessoas_sem_reclamacao;" +
		"")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize all variables
	var numIndicacoes int64 = 0
	var numPessoas int64 = 0
	var numConversas int64 = 0
	var numReprovados int64 = 0
	var indicacoesAprovadas int64 = 0
	var requerimentosAprovados int64 = 0
	var oficiosAprovados int64 = 0
	var indicacoes int64 = 0
	var numRequerimentos int64 = 0
	var numOficios int64 = 0
	var protocoladas int64 = 0
	var regioes []Regiao

	// Get counts from first query
	if rows.Next() {
		if err := rows.Scan(&numIndicacoes, &numConversas); err != nil {
			rows.Close()
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	rows.Close() // Close the first result set
	numPessoas = numIndicacoes + numConversas

	// Second query for approved items and protocols
	rows, err = db.Query("SELECT " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'indicação') AS indicacoes_aprovadas, " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'requerimento') AS requerimentos_aprovados, " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE status = 'aprovado' AND tipo = 'ofício') AS oficios_aprovados, " +
		"(SELECT COUNT(idReclamacao) FROM Protocolo) AS protocoladas;")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if rows.Next() {
		if err := rows.Scan(&indicacoesAprovadas, &requerimentosAprovados, &oficiosAprovados, &protocoladas); err != nil {
			rows.Close()
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	rows.Close()

	// Process region data
	for i := 1; i <= 6; i++ {
		mapeamento := make(map[string]int)
		for _, cat := range categorias {
			catRows, err := db.Query("SELECT COUNT(id) FROM Reclamacao WHERE categoria = $1 AND Regiao = $2", cat, i)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			contador := 0
			if catRows.Next() {
				if err := catRows.Scan(&contador); err != nil {
					catRows.Close()
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
			}
			catRows.Close()
			mapeamento[cat] = contador
		}
		regiao := Regiao{i, mapeamento}
		regioes = append(regioes, regiao)
	}

	// Get complaint types
	rows, err = db.Query("SELECT " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'indicação') AS numIndicacao, " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'ofício') AS numOficios, " +
		"(SELECT COUNT(id) FROM Reclamacao WHERE tipo = 'requerimento') AS numRequerimentos;")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rows.Next() {
		if err = rows.Scan(&indicacoes, &numOficios, &numRequerimentos); err != nil {
			rows.Close()
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	rows.Close()

	// Get rejected complaints
	rows, err = db.Query("SELECT COUNT(id) FROM Reclamacao WHERE status = 'reprovado';")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rows.Next() {
		if err = rows.Scan(&numReprovados); err != nil {
			rows.Close()
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	rows.Close()

	// Safe division to prevent divide by zero errors
	percIndicacao := 0.0
	if numPessoas > 0 {
		percIndicacao = float64(numIndicacoes) / float64(numPessoas) * 100
	}

	percProtocolacao := 0.0
	denominator := indicacoesAprovadas + requerimentosAprovados
	if denominator > 0 {
		percProtocolacao = float64(protocoladas) / float64(denominator) * 100
	}

	// Create stats
	stat := Stat{
		IndicacoesAprovadas:    int(indicacoesAprovadas),
		RequerimentosAprovados: int(requerimentosAprovados),
		OficiosAprovados:       int(oficiosAprovados),
		NumPessoas:             int(numPessoas),
		Indicacoes:             int(indicacoes),
		NumRequerimentos:       int(numRequerimentos),
		NumOficios:             int(numOficios),
		PercIndicacao:          percIndicacao,
		NumReclamacoes:         int(indicacoes + numRequerimentos),
		NumIndicacoes:          int(numIndicacoes),
		NumReprovados:          int(numReprovados),
		Protocoladas:           int(protocoladas),
		PercProtocolacao:       percProtocolacao,
		Regioes:                regioes,
	}

	c.IndentedJSON(http.StatusOK, stat)
}

func enviarNotificacao(telefone, nome, categoria string, numero int64) error {
	// Monta o conteúdo da mensagem
	msg := map[string]string{
		"telefone_cliente": telefone,
		"content": fmt.Sprintf(
			"Olá, %s :D! Só passando pra avisar que a sua solicitação de categoria %s foi protocolada pela Prefeitura Municipal de São Roque com o número %d.\nEspero que a situação seja resolvida o mais rápido possível, abraços!",
			nome, categoria, numero,
		),
	}

	// Serializa para JSON
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("erro ao gerar JSON: %w", err)
	}

	// Cria cliente HTTP com timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Cria a requisição
	req, err := http.NewRequest("POST", "http://147.93.10.167:5678/webhook/5f9d29a2-08bf-4dea-b7d2-9125eeaf48d9ARRUDAprotocolo", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Envia
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	// Loga resposta
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("Resposta do webhook (status %d): %s", resp.StatusCode, string(respBody))

	// Verifica status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("requisição falhou com status %d", resp.StatusCode)
	}

	return nil
}

func enviarProtocolo(c *gin.Context) {
	var protocolo Protocolo

	// Tenta converter o JSON recebido em um struct Protocolo
	if err := c.BindJSON(&protocolo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	// Insere o novo protocolo no banco de dados
	_, err := db.Exec("INSERT INTO protocolo (numero, idreclamacao) VALUES ($1, $2)", protocolo.Numero, protocolo.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Marca a reclamação como resolvida
	_, err = db.Exec("UPDATE reclamacao SET resolvido = TRUE WHERE id = $1", protocolo.ID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Busca dados da reclamação (nome, categoria, telefone) para notificação
	var nome, categoria, telefone string
	row := db.QueryRow("SELECT nome, categoria, telefone FROM reclamacao WHERE id = $1", protocolo.ID)
	if err := row.Scan(&nome, &categoria, &telefone); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Erro ao buscar dados da reclamação: " + err.Error()})
		return
	}

	// Envia a notificação via webhook N8N
	if err := enviarNotificacao(telefone, nome, categoria, protocolo.Numero); err != nil {
		log.Printf("Erro ao enviar notificação: %v", err)
	}

	// Marca o protocolo como avisado
	_, err = db.Exec("UPDATE protocolo SET avisado = TRUE WHERE idreclamacao = $1", protocolo.ID)
	if err != nil {
		log.Printf("Erro ao marcar protocolo como avisado: %v", err)
	}

	// Retorna resposta de sucesso
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Protocolo enviado e cliente notificado com sucesso."})
}

func login(c *gin.Context) {
	var usuario Usuario
	if err := c.BindJSON(&usuario); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}
	// Busca o usuário no banco de dados
	rows, err := db.Query("SELECT * FROM usuarios WHERE usuario = $1 AND senha = $2", usuario.Usuario, usuario.Senha)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	if rows.Next() {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Login sucesso.", "expiration_time": time.Now().Add(24 * time.Hour).Unix()})
		return
	}
	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Credenciais inválidas."})
}

func cadastrarEnderecos(c *gin.Context) {
	var data []Endereco
	if err := c.ShouldBindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Formato inválido.", "error": err.Error()})
		return
	}

	for _, endereco := range data {
		_, err := db.Exec("INSERT INTO enderecos (logradouro, bairro, regiao) VALUES ($1, $2, $3)", endereco.Logradouro, endereco.Bairro, endereco.Regiao)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Endereço cadastrado com sucesso!"})
}

func convertIndicacao(indicacao string) (string, error) {
	// URL base do webhook
	baseURL := "http://147.93.10.167:5678/webhook/a55c010b-4d45-47d4-a2b4-0bb8c5d60892indreqArruda"

	// Codificar a indicação corretamente
	indicacaoEncoded := url.QueryEscape(indicacao)

	// Agora construindo a URL
	webhookURL := fmt.Sprintf("%s?indicacao=%s", baseURL, indicacaoEncoded)

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	// Criar requisição POST sem corpo
	req, err := http.NewRequest("POST", webhookURL, nil)
	if err != nil {
		return "", err
	}

	// Enviar requisição
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Deserializar a resposta JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Message.Content, nil
}

func convertRequerimento(requerimento string) (string, error) {
	// URL base do webhook
	baseURL := "http://147.93.10.167:5678/webhook/a55c010b-4d45-47d4-a2b4-0bb8c5d60892INDICREQ2Arruda"

	// Codificar a indicação corretamente
	requerimentoEncoded := url.QueryEscape(requerimento)

	// Agora construindo a URL
	webhookURL := fmt.Sprintf("%s?requerimento=%s", baseURL, requerimentoEncoded)

	// Criar cliente HTTP com timeout
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	// Criar requisição POST sem corpo
	req, err := http.NewRequest("POST", webhookURL, nil)
	if err != nil {
		return "", err
	}

	// Enviar requisição
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Deserializar a resposta JSON
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Message.Content, nil
}

//leads ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

func getLeads(c *gin.Context) {
	var leads []Lead
	var ativo bool
	rows, err := db.Query("SELECT nome, telefone, ativo FROM Contatos")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var telefone, nome string

	for rows.Next() {
		if err := rows.Scan(&nome, &telefone, &ativo); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		leads = append(leads, Lead{
			Nome:     nome,
			Telefone: telefone,
			Ativo:    ativo,
		})
	}

	c.IndentedJSON(http.StatusOK, leads)
}

func getAllLeads(c *gin.Context) {
	var leads []Lead
	var ativo bool
	var count int
	rows, err := db.Query("SELECT nome, telefone, ativo FROM Contatos")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var telefone, nome string
	for rows.Next() {
		if err := rows.Scan(&nome, &telefone, &ativo); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if ativo {
			count++
		}
		leads = append(leads, Lead{
			Nome:     nome,
			Telefone: telefone,
			Ativo:    ativo,
		})
	}
	strCount := strconv.Itoa(count)
	total := strconv.Itoa(planos[plano_atual])
	ocupacao := strCount + "/" + total
	c.IndentedJSON(http.StatusOK, gin.H{"leads": leads, "ocupacao": ocupacao})
}

func desativarLead(c *gin.Context) {
	telefone := c.Param("telefone")
	_, err := db.Exec("UPDATE Contatos SET ativo = false WHERE telefone = $1", telefone)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.String(http.StatusOK, "Lead desativado com sucesso.")
}

func padronizaTelefone(telefone string) string {
	if len(telefone) > 13 && !strings.HasSuffix(telefone, "@lid") {
		telefone += "@lid"
	} else if !strings.HasSuffix(telefone, "@s.whatsapp.net") && !strings.HasSuffix(telefone, "@lid") {
		telefone += "@s.whatsapp.net"
	}
	return telefone
}

// cors ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	var dbError error

	// Configuração de conexão com o banco
	connStr := "host=localhost port=5432 user=evolution password=Albino29$$$$ dbname=clara_copia sslmode=disable"
	db, dbError = sql.Open("postgres", connStr)
	if dbError != nil {
		log.Fatal(dbError)
	}

	pingError := db.Ping()
	if pingError != nil {
		log.Fatal(pingError)
	}

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/send", messageHandler)
	r.PUT("/edit/:id", editarReclamacao)
	r.POST("/aprovar/:id", aprovarInquerito)
	r.POST("/aprovar/requerimento/:id", aprovarRequerimento)
	r.POST("/indicreq/:id", aprovarComoAmbos)
	r.POST("/reprovar/:id", reprovarInquerito)
	r.GET("/reclamacoes", getAllReclamacoes)
	r.GET("/stats", getStats)
	r.GET("/getRegiao", getRegiao)
	r.PUT("/protocolo", enviarProtocolo)
	r.POST("/login", login)
	r.GET("/", greetings)
	r.POST("/init", initializeContact)
	r.POST("/addMessage", addMessage)
	r.PUT("/setConversation", setConversation)
	r.GET("/getMessages/:telefone", getMessagesByPhone)
	r.DELETE("/clearMessages/:telefone", clearMessagesByPhone)
	r.POST("/cadastrarEnderecos", cadastrarEnderecos)
	r.POST("/aprovar/oficio/:id", aprovarOficio)
	r.GET("/getLeads", getLeads)
	r.GET("/getAllLeads", getAllLeads)
	r.DELETE("/desativarLead/:telefone", desativarLead)
	fmt.Println("Servidor rodando em http://localhost:3017")
	log.Fatal(r.Run(":3017")) // Inicia o servidor na porta 3017
}
