package services

import (
	"back-end/config"
	"back-end/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Responder(telefone, conversationID, mensagem string) error {
	baseURL := os.Getenv("WEBHOOK_RESPONDER")
	client := &http.Client{Timeout: time.Second * 10}

	telefonePadronizado := config.PadronizaTelefone(telefone)

	data := url.Values{}
	data.Set("user", telefonePadronizado)
	data.Set("data", mensagem)
	data.Set("telefone", telefonePadronizado)
	data.Set("conversation_id", conversationID)

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

func converter(baseURL, param, valor string) (string, error) {
	webhookURL := fmt.Sprintf("%s?%s=%s", baseURL, param, url.QueryEscape(valor))
	client := &http.Client{Timeout: time.Second * 60}

	req, err := http.NewRequest("POST", webhookURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response models.WebhookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.Message.Content, nil
}

func ConvertIndicacao(indicacao string) (string, error) {
	return converter(os.Getenv("WEBHOOK_INDICACAO"), "indicacao", indicacao)
}

func ConvertRequerimento(requerimento string) (string, error) {
	return converter(os.Getenv("WEBHOOK_REQUERIMENTO"), "requerimento", requerimento)
}

func EnviarNotificacao(telefone, nome, categoria string, numero int64) error {
	mensagem := fmt.Sprintf(
		"Olá, %s :D! Só passando pra avisar que a sua solicitação de categoria %s foi protocolada pela Prefeitura Municipal de Sorocaba com o número %d.\nEspero que a situação seja resolvida o mais rápido possível, abraços!",
		nome, categoria, numero,
	)

	baseURL := os.Getenv("WEBHOOK_ENVIAR_MENSAGEM")
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	data.Set("mensagem", mensagem)
	data.Set("telefone", telefone)
	data.Set("instance", "ouvidoria_clara_jussara")

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
