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
func (c *Config) Client() (*mailgun.Mailgun, error) {

	domain := "" // We don't set a domain right away
	apiKey := c.APIKey
	publicAPIKey := "" // We don't support email validation

	client := mailgun.NewMailgun(domain, apiKey, publicAPIKey)

	// if err != nil {
	// 	return nil, err
	// }

	log.Printf("[INFO] Mailgun Client configured ")

	return &client, nil
}
