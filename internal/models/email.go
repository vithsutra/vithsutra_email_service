package models

type Email struct {
	To          string            `json:"to"`
	Subject     string            `json:"subject"`
	ServiceName string            `json:"service_name"`
	EmailType   string            `json:"email_type"`
	Data        map[string]string `json:"data"`
}
