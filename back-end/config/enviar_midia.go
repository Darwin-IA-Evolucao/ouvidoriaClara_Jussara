package config

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func EnviarMidia(telefone, mensagem string) error {
	baseURL := os.Getenv("WEBHOOK_ENVIAR_MENSAGEM")
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	data := url.Values{}
	data.Set("mensagem", mensagem)
	data.Set("number", telefone)
	data.Set("mediatype","")
	data.Set("media","")
	data.Set("filename","")
	data.Set("caption","")
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