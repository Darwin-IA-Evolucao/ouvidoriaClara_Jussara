package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type PlanilhaLinha struct {
	Linha    int
	Nome     string
	Telefone string
}

var telefoneRegex = regexp.MustCompile(`^55\d{10,11}$`)

func ParsePlanilhaContatos(filename string, r io.Reader) ([]PlanilhaLinha, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo")
	}
	switch ext {
	case ".csv":
		return parseCSVContatos(data)
	case ".xlsx":
		return parseXLSXContatos(data)
	default:
		return nil, fmt.Errorf("tipo de arquivo nao suportado, envie csv ou xlsx")
	}
}

func NormalizeTelefone(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", "+", "", "\u00a0", "")
	s = replacer.Replace(s)
	if strings.ContainsAny(s, ".eE") {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			s = strconv.FormatInt(int64(f), 10)
		}
	}
	return s
}

func TelefoneValido(s string) bool {
	return telefoneRegex.MatchString(s)
}

func parseCSVContatos(data []byte) ([]PlanilhaLinha, error) {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	firstLine, _, _ := bytes.Cut(data, []byte("\n"))
	firstLine = bytes.TrimSuffix(firstLine, []byte("\r"))
	delim := ','
	if bytes.Count(firstLine, []byte(";")) > bytes.Count(firstLine, []byte(",")) {
		delim = ';'
	}
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = delim
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csv invalido")
	}
	return mapPlanilhaRows(records)
}

func parseXLSXContatos(data []byte) ([]PlanilhaLinha, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("xlsx invalido")
	}
	defer f.Close()
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("xlsx sem planilhas")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("erro ao ler xlsx")
	}
	return mapPlanilhaRows(rows)
}

func mapPlanilhaRows(rows [][]string) ([]PlanilhaLinha, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("planilha vazia")
	}
	nomeIdx, telIdx := -1, -1
	for i, h := range rows[0] {
		switch normHeader(h) {
		case "nome":
			nomeIdx = i
		case "telefone":
			telIdx = i
		}
	}
	if nomeIdx < 0 || telIdx < 0 {
		return nil, fmt.Errorf("planilha deve ter colunas nome e telefone")
	}
	out := make([]PlanilhaLinha, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		nome := cellAt(rows[i], nomeIdx)
		tel := cellAt(rows[i], telIdx)
		if nome == "" && tel == "" {
			continue
		}
		out = append(out, PlanilhaLinha{Linha: i + 1, Nome: nome, Telefone: tel})
	}
	return out, nil
}

func normHeader(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.TrimPrefix(s, "\ufeff")
	return s
}

func cellAt(row []string, i int) string {
	if i < 0 || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}
