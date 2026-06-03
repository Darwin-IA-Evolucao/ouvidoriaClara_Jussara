package models

type WebhookMensagem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type WebhookResponse struct {
	Index        int             `json:"index"`
	Message      WebhookMensagem `json:"message"`
	FinishReason string          `json:"finish_reason"`
}
