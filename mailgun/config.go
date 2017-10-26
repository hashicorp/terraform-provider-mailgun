package mailgun

import (
	"log"

	"github.com/mailgun/mailgun-go"
)

// Config is the configuration of mailgun provider
type Config struct {
	APIKey string
}

// Client returns a new client for accessing mailgun.
func (c *Config) Client() (mailgun.Mailgun, error) {
	client := mailgun.NewMailgun("", c.APIKey, "")

	log.Printf("[INFO] Mailgun Client configured ")

	return client, nil
}
