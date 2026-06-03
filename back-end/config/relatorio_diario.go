package config

import (
	"fmt"
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
		(SELECT COUNT(*) FROM interesse WHERE datainteresse >= now() - INTERVAL '1 day') AS total_interesse,
		(SELECT COUNT(*) FROM contatos WHERE data_criacao >= now() - INTERVAL '24 hours' AND (ehservidorouoficial = true OR ehservidorouoficial is null)) AS total_contatos,
		(SELECT COUNT(DISTINCT c.telefone)
		FROM contatos c
		WHERE c.data_criacao >= now() - INTERVAL '24 hours'
			AND (c.ehservidorouoficial = true OR c.ehservidorouoficial is null)
			AND (  
					EXISTS (
						SELECT 1 FROM interesse i
						WHERE i.fk_cliente_telefonecliente = c.telefone
						AND i.datainteresse >= now() - INTERVAL '24 hours'
					)
				)
		) AS total_contatos_convertidos,
		(SELECT COUNT(DISTINCT c.telefone)
			FROM contatos c
			LEFT JOIN interesse i
			ON c.telefone = i.fk_cliente_telefonecliente
			AND i.datainteresse >= now() - INTERVAL '24 hours'
			WHERE c.data_criacao >= now() - INTERVAL '24 hours'
			AND i.fk_cliente_telefonecliente IS NULL
			AND (c.ehservidorouoficial = true OR c.ehservidorouoficial is null)
		) AS total_contatos_reativar,
		(SELECT COUNT(*) FROM cliente WHERE data_criacao >= now() - INTERVAL '24 hours') AS total_clientes_novos,
		(SELECT COUNT(DISTINCT cl.telefonecliente)
		FROM cliente cl
		WHERE cl.data_criacao < now() - INTERVAL '24 hours'
			AND (
					EXISTS (
						SELECT 1 FROM interesse i
						WHERE i.fk_cliente_telefonecliente = cl.telefonecliente
						AND i.datainteresse >= now() - INTERVAL '24 hours'
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

	msg := "📊 *Relatório Diário Darwin IA* 📊\n\n"
	msg += fmt.Sprintf("Olá ! Aqui é a Laura sua IA, segue o relatório de hoje (%s):\n\n", hoje)
	msg += "🔥 *Interesses (24h)*\n"
	msg += fmt.Sprintf("- Total: %d\n\n", dados.TotalInteresse)
	msg += "📞 *Contatos (24h)*\n"
	msg += fmt.Sprintf("- Recebidos: %d\n", dados.TotalContatos)
	msg += fmt.Sprintf("- Convertidos em interesse: %d\n", dados.TotalConvertidos)
	msg += fmt.Sprintf("- Para reativação: %d\n\n", dados.TotalReativar)
	msg += "👥 *Clientes*\n"
	msg += fmt.Sprintf("- Novos cadastros (24h): %d\n", dados.TotalClientesNovos)
	msg += fmt.Sprintf("- Cliente antigo: %d", dados.TotalClientesRecorrentes)

	return msg
}

func RelatorioDiario(conn *sqlx.DB) {
	telefone := "5515998223027"
	for {
		agora := time.Now()
		proxima := time.Date(agora.Year(), agora.Month(), agora.Day(), 10, 18, 0, 0, agora.Location())
		if !agora.Before(proxima) {
			proxima = proxima.Add(24 * time.Hour)
		}
		time.Sleep(time.Until(proxima))
		dados, err := GetRelatorioDiario(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if err := EnviarRelatorio(telefone, GerarMensagemRelatorio(*dados)); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func GetRelatorioGelando(conn *sqlx.DB) (*Relatorio, error) {
	const query = `SELECT count(*) AS total_contatos
FROM contatos c
LEFT JOIN interesse cl
ON c.telefone = cl.fk_cliente_telefonecliente
WHERE cl.fk_cliente_telefonecliente IS NULL
AND data_criacao > CURRENT_DATE - INTERVAL '24 hours'
AND (c.ehservidorouoficial = true OR c.ehservidorouoficial IS NULL);`
	var relatorio Relatorio
	err := conn.Get(&relatorio, query)
	if err != nil {
		return nil, err
	}

	return &relatorio, nil
}

func GerarMensagemGelando(dados Relatorio) string {
	hoje := time.Now().Format("02/01/2006")

	msg := "📊🧊 *Relatório Clientes Gelando Darwin IA* 🧊📊\n\n"
	msg += fmt.Sprintf("Olá ! Aqui é a Laura sua IA, notamos que há clientes gelando no dia de hoje (%s):\n\n", hoje)

	msg += fmt.Sprintf("- Total de clientes gelando: %d\n", dados.TotalContatos)

	msg += "Sugerimos entrar no painel de controle e verificar os clientes gelando para reativar o contato.\n\n"
	msg += "Link do painel: https://unique.sdr.darwinsistema.com.br"

	return msg
}

func RelatorioGelando(conn *sqlx.DB) {
	telefone1 := "5515998223027"
	telefone2 := "5515998223027"
	for {
		agora := time.Now()
		proxima := time.Date(agora.Year(), agora.Month(), agora.Day(), 10, 16, 0, 0, agora.Location())
		if !agora.Before(proxima) {
			proxima = proxima.Add(24 * time.Hour)
		}
		time.Sleep(time.Until(proxima))
		dados, err := GetRelatorioGelando(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if err := EnviarRelatorio(telefone1, GerarMensagemGelando(*dados)); err != nil {
			fmt.Println(err.Error())
		}
		if err := EnviarRelatorio(telefone2, GerarMensagemGelando(*dados)); err != nil {
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
    (SELECT COUNT(*) FROM interesse i, periodo p WHERE datainteresse >= p.inicio_mes_passado AND datainteresse < p.inicio_mes_atual) AS total_interesse,
    (SELECT COUNT(*) FROM contatos c, periodo p WHERE data_criacao >= p.inicio_mes_passado AND data_criacao < p.inicio_mes_atual AND (c.ehservidorouoficial = true OR c.ehservidorouoficial is null)) AS total_contatos,
    (SELECT COUNT(DISTINCT c.telefone) FROM contatos c, periodo p 
        WHERE c.data_criacao >= p.inicio_mes_passado AND c.data_criacao < p.inicio_mes_atual
        AND (c.ehservidorouoficial = true OR c.ehservidorouoficial is null)
        AND (  
			EXISTS (
			    SELECT 1 FROM interesse i
				WHERE i.fk_cliente_telefonecliente = c.telefone
				AND i.datainteresse >= p.inicio_mes_passado
				AND i.datainteresse < p.inicio_mes_atual
			)
		)
    ) AS total_contatos_convertidos,
    (SELECT COUNT(DISTINCT c.telefone)
		FROM contatos c
        CROSS JOIN periodo p
		LEFT JOIN interesse i
        ON c.telefone = i.fk_cliente_telefonecliente    
        AND i.datainteresse >= p.inicio_mes_passado
        AND i.datainteresse < p.inicio_mes_atual
		WHERE c.data_criacao >= p.inicio_mes_passado 
        AND c.data_criacao < p.inicio_mes_atual   
        AND (c.ehservidorouoficial = true OR c.ehservidorouoficial is null)
	    AND i.fk_cliente_telefonecliente IS NULL) AS total_contatos_reativar,
    (SELECT COUNT(*) FROM cliente cl, periodo p WHERE data_criacao >= p.inicio_mes_passado AND data_criacao < p.inicio_mes_atual) AS total_clientes_novos,
    (SELECT COUNT(DISTINCT cl.telefonecliente)
		FROM cliente cl, periodo p
		WHERE data_criacao < p.inicio_mes_passado
			AND (
				EXISTS (
					SELECT 1 FROM interesse i
					WHERE i.fk_cliente_telefonecliente = cl.telefonecliente
					AND i.datainteresse >= p.inicio_mes_passado
					AND i.datainteresse < p.inicio_mes_atual
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

	msg := "📊 Relatório Mensal Darwin IA 📊\n\n"
	msg += fmt.Sprintf("Olá ! Aqui é a Laura, sua IA, segue o relatório do mês (%s):\n\n", hoje)
	msg += "🔥 Interesses\n"
	msg += fmt.Sprintf("- Total: %d\n\n", dados.TotalInteresse)
	msg += "📞 Contatos\n"
	msg += fmt.Sprintf("- Recebidos: %d\n", dados.TotalContatos)
	msg += fmt.Sprintf("- Convertidos em interesse: %d\n", dados.TotalConvertidos)
	msg += fmt.Sprintf("- Para reativação: %d\n\n", dados.TotalReativar)
	msg += "👥 Clientes\n"
	msg += fmt.Sprintf("- Novos cadastros: %d\n", dados.TotalClientesNovos)
	msg += fmt.Sprintf("- Cliente antigo: %d", dados.TotalClientesRecorrentes)

	return msg
}

func RelatorioMensal(conn *sqlx.DB) {
	telefone := "5515998223027"
	for {
		agora := time.Now()
		proximoMes := agora.Month() + 1
		if proximoMes > 12 {
			proximoMes = time.January
			agora = agora.AddDate(1, 0, 0) // Avança para o próximo ano
		}

		proxima := time.Date(agora.Year(), proximoMes, 1, 10, 17, 0, 0, agora.Location())
		fmt.Printf("[RELATORIO MENSAL] Próximo envio programado para: %s\n", proxima.Format("02/01/2006 15:04:05"))
		time.Sleep(time.Until(proxima))

		dados, err := GetRelatorioMensal(conn)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if err := EnviarRelatorio(telefone, GerarMensagemRelatorioMensal(*dados)); err != nil {
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
