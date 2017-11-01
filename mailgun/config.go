package mailgun

import (
	"log"

	"gopkg.in/mailgun/mailgun-go.v1"
)

// Config - configuration struct for mailgun
type Config struct {
	APIKey string
}

// Client() returns a new client for accessing mailgun.

func (c *Config) Client() (*MailgunProvider, error) {
	client := MailgunProvider{
		Config:  c,
		clients: make(map[string]mailgun.Mailgun),
	}

	// if err != nil {
	// 	return nil, err
	// }

	log.Printf("[INFO] Mailgun Client configured ")

	return &client, nil
}

type MailgunProvider struct {
	*Config
	clients map[string]mailgun.Mailgun
}

func (p *MailgunProvider) Domain(domain string) mailgun.Mailgun {
	publicAPIKey := "" // We don't support email validation
	var client mailgun.Mailgun
	var ok bool
	if client, ok = p.clients[domain]; !ok {
		client = mailgun.NewMailgun(domain, p.APIKey, publicAPIKey)
		p.clients[domain] = client
	}
	return client
}
