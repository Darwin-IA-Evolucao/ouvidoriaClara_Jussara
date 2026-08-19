package config

import (
	"back-end/models"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

func VerificarLembrete(conn *sqlx.DB) ([]models.Cliente, error) {
	const query = `SELECT c.telefone,  c.nome FROM cliente c 
				JOIN reclamacao r ON r.telefone = c.telefone
				JOIN aviso a ON a.id_reclamacao = r.idreclamacao
				WHERE a.avisado = false AND a.data BETWEEN NOW() - INTERVAL '1 day' AND NOW()`

	var clientes []models.Cliente
	err := conn.Select(&clientes, query)
	return clientes, err
}

func AtualizarAviso(conn *sqlx.DB) error {
	const query = `UPDATE aviso SET avisado = true WHERE data BETWEEN NOW() - INTERVAL '1 day' AND NOW() AND avisado = false`
	_, err := conn.Exec(query)
	return err
}
func EnviarLembrete(conn *sqlx.DB) {
	clientes, err := VerificarLembrete(conn)
	if err != nil {
		fmt.Println("[ENVIAR LEMBRETE] Erro ao verificar lembrete: ", err)
		return
	}
	fmt.Printf("[ENVIAR LEMBRETE] Encontrados %d clientes com lembrete para enviar\n", len(clientes))
	// Enviar lembrete para cada cliente
	for _, cliente := range clientes {
		fmt.Printf("[ENVIAR LEMBRETE] Enviando lembrete para: %s - %s\n", cliente.Telefone, cliente.Nome)
		msg := "🤝LEMBRETE🤝\n"
		msg += fmt.Sprintf("Olá! O(a) cliente %s tem uma solicitação com prazo marcado para hoje.\n", cliente.Nome)
		msg += fmt.Sprintf("Entre em contato com ele(a) pelo número: %s.", cliente.Telefone)
		err = EnviarMensagem(os.Getenv("TELEFONE_GERAL"), msg)
		if err != nil {
			fmt.Println("[ENVIAR LEMBRETE] Erro ao enviar lembrete: ", err)
			continue
		}
		fmt.Printf("[ENVIAR LEMBRETE] Lembrete enviado para: %s - %s\n", cliente.Telefone, cliente.Nome)
	}
	err = AtualizarAviso(conn)
	if err != nil {
		fmt.Println("[ENVIAR LEMBRETE] Erro ao atualizar aviso: ", err)
		return
	}
}

func EnviarLembreteNoHorario(conn *sqlx.DB) {
	for {
		now := time.Now()
		prox := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

		if now.Before(prox) {
			fmt.Println("[ENVIAR LEMBRETE] Aguardando até:", prox)
			time.Sleep(time.Until(prox))
		}

		EnviarLembrete(conn)

		prox = prox.Add(24 * time.Hour)
		fmt.Println("[ENVIAR LEMBRETE] Aguardando até:", prox)
		time.Sleep(time.Until(prox))
	}
}
