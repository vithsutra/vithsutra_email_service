package models

type Email struct {
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	EmailType string            `json:"email_type"`
	Data      map[string]string `json:"data"`
}
