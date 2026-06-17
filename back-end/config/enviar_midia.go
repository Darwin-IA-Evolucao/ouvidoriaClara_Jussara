package config

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func EnviarMidia(telefone, mensagem, link string) error {
	baseURL := os.Getenv("WEBHOOK_ENVIAR_MENSAGEM")
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	mimeType, mediaType, fileName := separarLink(link)
	data := url.Values{}
	data.Set("caption", mensagem)
	data.Set("number", telefone)
	data.Set("mediatype", mediaType)
	data.Set("media", link)
	data.Set("filename", fileName)
	data.Set("mimetype", mimeType)

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

func separarLink(link string) (mimetype string, mediatype string, fileName string){
	linkSplit := strings.TrimPrefix(filepath.Ext(link), ".")
	fileName = strings.TrimSuffix(filepath.Base(link), filepath.Ext(link))
	switch linkSplit {
	case "png", "jpg", "jpeg":
		mimetype = "image/png"
		mediatype = "image"
	case "mp4", "avi", "mov":
		mimetype = "video/mp4"
		mediatype = "video"
	case "pdf":
		mimetype = "application/pdf"
		mediatype = "document"
	}
	fileName = fileName+"."+linkSplit
	return mimetype, mediatype, fileName
}
