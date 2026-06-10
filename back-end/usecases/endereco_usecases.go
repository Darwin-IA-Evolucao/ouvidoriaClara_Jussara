package usecases

import (
	"back-end/models"
	"back-end/repository"
	"strings"
)

type EnderecoUseCases struct {
	repository repository.EnderecoRepository
}

func NewEnderecoUseCases(repo repository.EnderecoRepository) EnderecoUseCases {
	return EnderecoUseCases{repository: repo}
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
	input = strings.Split(input, ",")[0]
	return strings.TrimSpace(input)
}

func buscaBinaria(lista []models.Logradouro, ruaProc string) (int, bool) {
	var ruaAtual string
	ruaProc = strings.ToLower(ruaProc)
	inicio := 0
	fim := len(lista) - 1
	meio := (fim + inicio) / 2

	for fim >= inicio {
		ruaAtual = strings.ToLower(lista[meio].Logradouro)
		if ruaAtual == ruaProc {
			return meio, true
		}
		if ruaAtual > ruaProc {
			fim = meio - 1
			meio = inicio + (fim-inicio)/2
		} else {
			inicio = meio + 1
			meio = inicio + (fim-inicio)/2
		}
	}
	return -1, false
}

func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

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
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min3(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[rows-1][cols-1]
}

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

func encontraMelhorCorrespondencia(lista []models.Logradouro, termo string, limiteSimilaridade int) (int, bool) {
	termoNormalizado := strings.ToLower(strings.TrimSpace(termo))
	if termoNormalizado == "" {
		return -1, false
	}

	if limiteSimilaridade <= 0 {
		limiteSimilaridade = len(termo)
		if limiteSimilaridade > 10 {
			limiteSimilaridade = 10
		}
	}

	melhorIndice := -1
	melhorPontuacao := -1

	for i, item := range lista {
		logradouroNormalizado := strings.ToLower(strings.TrimSpace(item.Logradouro))
		distancia := levenshteinDistance(termoNormalizado, logradouroNormalizado)
		pontuacao := distancia

		if strings.Contains(logradouroNormalizado, termoNormalizado) {
			pontuacao = pontuacao / 2
		}

		termoPalavras := strings.Fields(termoNormalizado)
		palavrasEncontradas := 0
		for _, palavra := range termoPalavras {
			if strings.Contains(logradouroNormalizado, palavra) {
				palavrasEncontradas++
			}
		}

		if palavrasEncontradas == len(termoPalavras) && len(termoPalavras) > 0 {
			pontuacao = pontuacao / 3
		} else if palavrasEncontradas > 0 {
			pontuacao = pontuacao * (len(termoPalavras) - palavrasEncontradas + 1) / len(termoPalavras)
		}

		if melhorPontuacao == -1 || pontuacao < melhorPontuacao {
			melhorPontuacao = pontuacao
			melhorIndice = i
		}
	}

	if melhorPontuacao <= limiteSimilaridade || melhorPontuacao < 5 {
		return melhorIndice, true
	}
	return -1, false
}

func (uc EnderecoUseCases) 	GetRegiao(input string) int {
	input = limparEndereco(input)

	if regiao, err := uc.repository.GetRegiaoByLogradouro(input); err == nil {
		return regiao
	}

	logradouros, err := uc.repository.GetAllLogradouros()
	if err != nil {
		return 1
	}

	if index, achou := buscaBinaria(logradouros, input); achou {
		return logradouros[index].Regiao
	}

	limiteSimilaridade := len(input)
	if limiteSimilaridade > 10 {
		limiteSimilaridade = 10
	}
	if melhorIndex, encontrado := encontraMelhorCorrespondencia(logradouros, input, limiteSimilaridade); encontrado {
		return logradouros[melhorIndex].Regiao
	}

	return 1
}

func (uc EnderecoUseCases) CadastrarEnderecos(enderecos []models.Endereco) error {
	for _, endereco := range enderecos {
		if err := uc.repository.CreateEndereco(endereco); err != nil {
			return err
		}
	}
	return nil
}
