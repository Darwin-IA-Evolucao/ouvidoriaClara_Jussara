package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type Relatorio struct {
	TotalInteresse           int `db:"total_interesse"`
	TotalContatos            int `db:"total_contatos"`
	TotalConvertidos         int `db:"total_contatos_convertidos"`
	TotalReativar            int `db:"total_contatos_reativar"`
	TotalClientesNovos       int `db:"total_clientes_novos"`
	TotalClientesRecorrentes int `db:"total_clientes_recorrentes"`
}

func GetRelatorioDiario(conn *sqlx.DB) (*Relatorio, error) {
	const query = `
	SELECT
		(SELECT COUNT(*) FROM reclamacao WHERE data_criacao >= now() - INTERVAL '1 day') AS total_interesse,

		(SELECT COUNT(*) FROM contatos WHERE data_criacao >= now() - INTERVAL '24 hours') AS total_contatos,

		(
			SELECT COUNT(DISTINCT c.telefone)
			FROM contatos c
			WHERE c.data_criacao >= now() - INTERVAL '24 hours'
			AND (
				EXISTS (
					SELECT 1
					FROM reclamacao i
					WHERE i.telefone = c.telefone
					AND i.data_criacao >= now() - INTERVAL '24 hours'
				)
			)
		) AS total_contatos_convertidos,

		(
			SELECT COUNT(DISTINCT c.telefone)
			FROM contatos c

			LEFT JOIN reclamacao i
				ON c.telefone = i.telefone
				AND i.data_criacao >= now() - INTERVAL '24 hours'

			WHERE c.data_criacao >= now() - INTERVAL '24 hours'
			AND i.telefone IS NULL
		) AS total_contatos_reativar,

		(
			SELECT COUNT(*)
			FROM cliente
			WHERE data_criacao >= now() - INTERVAL '24 hours'
		) AS total_clientes_novos,

		(
			SELECT COUNT(DISTINCT cl.telefone)
			FROM cliente cl
			WHERE cl.data_criacao < now() - INTERVAL '24 hours'
			AND (
				EXISTS (
					SELECT 1
					FROM reclamacao i
					WHERE i.telefone = cl.telefone
					AND i.data_criacao >= now() - INTERVAL '24 hours'
				)
			)
		) AS total_clientes_recorrentes;
	`

	var relatorio Relatorio

	err := conn.Get(&relatorio, query)
	if err != nil {
		return nil, err
	}

	return &relatorio, nil
}

func GerarMensagemRelatorio(dados Relatorio) string {
	hoje := time.Now().Format("02/01/2006")

	msg := fmt.Sprintf(
		"📊 *Relatório Diário Darwin IA* 📊 • "+
			"Olá! Aqui é a Ju, sua IA. Segue o relatório de hoje (%s): • "+
			"🔥 *Interesses (24h)* • Total: %d • "+
			"📞 *Contatos (24h)* • Recebidos: %d • Convertidos em reclamação: %d • Para auxiliar na finalização: %d • "+
			"👥 *Clientes* • Novos cadastros (24h): %d • Clientes antigos: %d",
		hoje,
		dados.TotalInteresse,
		dados.TotalContatos,
		dados.TotalConvertidos,
		dados.TotalReativar,
		dados.TotalClientesNovos,
		dados.TotalClientesRecorrentes,
	)

	return msg
}

func RelatorioDiario(conn *sqlx.DB) {
	telefone := os.Getenv("TELEFONE_GERAL")

	for {
		agora := time.Now()

		proxima := time.Date(
			agora.Year(),
			agora.Month(),
			agora.Day(),
			18,
			0,
			0,
			0,
			agora.Location(),
		)

		if !agora.Before(proxima) {
			proxima = proxima.Add(24 * time.Hour)
		}

		fmt.Printf(
			"[RELATORIO DIARIO] Próximo envio programado para: %s\n",
			proxima.Format("02/01/2006 15:04:05"),
		)

		time.Sleep(time.Until(proxima))

		dados, err := GetRelatorioDiario(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if err := EnviarRelatorio(
			telefone,
			GerarMensagemRelatorio(*dados),
			"Relatório Diário",
		); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetRelatorioGelando(conn *sqlx.DB) (*Relatorio, error) {
	const query = `
	SELECT count(*) AS total_contatos
	FROM contatos c

	LEFT JOIN reclamacao cl
		ON c.telefone = cl.telefone

	WHERE cl.telefone IS NULL
	AND c.data_criacao > CURRENT_DATE - INTERVAL '24 hours';
	`

	var relatorio Relatorio

	err := conn.Get(&relatorio, query)
	if err != nil {
		return nil, err
	}

	return &relatorio, nil
}

func GerarMensagemGelando(dados Relatorio) string {
	hoje := time.Now().Format("02/01/2006")

	msg := fmt.Sprintf(
		"📊🧊 *Relatório Clientes Gelando Darwin IA* 🧊📊 • "+
			"Olá! Aqui é a Ju, sua IA. Notamos que há clientes que não finalizaram suas reclamações e estão virando clientes gelados hoje (%s). • "+
			"Total de clientes gelando: %d • "+
			"Sugerimos acessar o painel de controle e verificar os contatos que estão gelando para auxiliar na finalização. • "+
			"Link do painel: https://vereadorajussara.ouvidoria.darwinsistema.com.br",
		hoje,
		dados.TotalContatos,
	)

	return msg
}

func RelatorioGelando(conn *sqlx.DB) {
	telefone := os.Getenv("TELEFONE_GERAL")

	for {
		agora := time.Now()

		proxima := time.Date(
			agora.Year(),
			agora.Month(),
			agora.Day(),
			8,
			0,
			0,
			0,
			agora.Location(),
		)

		if !agora.Before(proxima) {
			proxima = proxima.Add(24 * time.Hour)
		}

		fmt.Printf(
			"[RELATORIO GELO] Próximo envio programado para: %s\n",
			proxima.Format("02/01/2006 15:04:05"),
		)

		time.Sleep(time.Until(proxima))

		dados, err := GetRelatorioGelando(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if err := EnviarRelatorio(
			telefone,
			GerarMensagemGelando(*dados),
			"Relatório Gelo",
		); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetRelatorioMensal(conn *sqlx.DB) (*Relatorio, error) {
	const query = `
	WITH periodo AS (
		SELECT
			date_trunc('month', now()) - interval '1 month' AS inicio_mes_passado,
			date_trunc('month', now()) AS inicio_mes_atual
	)

	SELECT

	(
		SELECT COUNT(*)
		FROM reclamacao i, periodo p
		WHERE data_criacao >= p.inicio_mes_passado
		AND data_criacao < p.inicio_mes_atual
	) AS total_interesse,

	(
		SELECT COUNT(*)
		FROM contatos c, periodo p
		WHERE data_criacao >= p.inicio_mes_passado
		AND data_criacao < p.inicio_mes_atual
	) AS total_contatos,

	(
		SELECT COUNT(DISTINCT c.telefone)
		FROM contatos c, periodo p

		WHERE c.data_criacao >= p.inicio_mes_passado
		AND c.data_criacao < p.inicio_mes_atual

		AND (
			EXISTS (
				SELECT 1
				FROM reclamacao i
				WHERE i.telefone = c.telefone
				AND i.data_criacao >= p.inicio_mes_passado
				AND i.data_criacao < p.inicio_mes_atual
			)
		)
	) AS total_contatos_convertidos,

	(
		SELECT COUNT(DISTINCT c.telefone)
		FROM contatos c

		CROSS JOIN periodo p

		LEFT JOIN reclamacao i
			ON c.telefone = i.telefone
			AND i.data_criacao >= p.inicio_mes_passado
			AND i.data_criacao < p.inicio_mes_atual

		WHERE c.data_criacao >= p.inicio_mes_passado
		AND c.data_criacao < p.inicio_mes_atual
		AND i.telefone IS NULL
	) AS total_contatos_reativar,

	(
		SELECT COUNT(*)
		FROM cliente cl, periodo p
		WHERE data_criacao >= p.inicio_mes_passado
		AND data_criacao < p.inicio_mes_atual
	) AS total_clientes_novos,

	(
		SELECT COUNT(DISTINCT cl.telefone)
		FROM cliente cl, periodo p

		WHERE data_criacao < p.inicio_mes_passado

		AND (
			EXISTS (
				SELECT 1
				FROM reclamacao i
				WHERE i.telefone = cl.telefone
				AND i.data_criacao >= p.inicio_mes_passado
				AND i.data_criacao < p.inicio_mes_atual
			)
		)
	) AS total_clientes_recorrentes;
	`

	var relatorio Relatorio

	err := conn.Get(&relatorio, query)
	if err != nil {
		return nil, err
	}

	return &relatorio, nil
}

func GerarMensagemRelatorioMensal(dados Relatorio) string {
	hoje := time.Now().Month().String()
	hoje = MonthTranslation(hoje)
	hoje += "/" + fmt.Sprint(time.Now().Year())

	msg := fmt.Sprintf(
		"📊 *Relatório Mensal Darwin IA* 📊 • "+
			"Olá! Aqui é a Ju, sua IA. Segue o relatório do mês (%s): • "+
			"🔥 *Interesses* • Total: %d • "+
			"📞 *Contatos* • Recebidos: %d • Convertidos em reclamação: %d • Para reativação: %d • "+
			"👥 *Clientes* • Novos cadastros: %d • Clientes antigos: %d",
		hoje,
		dados.TotalInteresse,
		dados.TotalContatos,
		dados.TotalConvertidos,
		dados.TotalReativar,
		dados.TotalClientesNovos,
		dados.TotalClientesRecorrentes,
	)

	return msg
}

func RelatorioMensal(conn *sqlx.DB) {
	telefone := os.Getenv("TELEFONE_GERAL")

	for {
		agora := time.Now()

		proxima := time.Date(
			agora.Year(),
			agora.Month(),
			1,
			18,
			0,
			0,
			0,
			agora.Location(),
		)

		// Se o dia 1 às 18h já passou,
		// agenda para o dia 1 do próximo mês.
		if !proxima.After(agora) {
			proxima = proxima.AddDate(0, 1, 0)
		}

		fmt.Printf(
			"[RELATORIO MENSAL] Próximo envio programado para: %s\n",
			proxima.Format("02/01/2006 15:04:05"),
		)

		time.Sleep(time.Until(proxima))

		dados, err := GetRelatorioMensal(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if err := EnviarRelatorio(
			telefone,
			GerarMensagemRelatorioMensal(*dados),
			"Relatório Mensal",
		); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func MonthTranslation(m string) string {
	switch m {
	case "January":
		return "Janeiro"

	case "February":
		return "Fevereiro"

	case "March":
		return "Março"

	case "April":
		return "Abril"

	case "May":
		return "Maio"

	case "June":
		return "Junho"

	case "July":
		return "Julho"

	case "August":
		return "Agosto"

	case "September":
		return "Setembro"

	case "October":
		return "Outubro"

	case "November":
		return "Novembro"

	case "December":
		return "Dezembro"
	}

	return ""
}
