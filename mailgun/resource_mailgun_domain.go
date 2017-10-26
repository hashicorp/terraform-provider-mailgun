package mailgun

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mailgun/mailgun-go"
)

func resourceMailgunDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceMailgunDomainCreate,
		Read:   resourceMailgunDomainRead,
		Delete: resourceMailgunDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spam_action": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},

			"smtp_password": &schema.Schema{
				Type:      schema.TypeString,
				ForceNew:  true,
				Required:  true,
				Sensitive: true,
			},

			"smtp_login": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"wildcard": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceMailgunDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(mailgun.Mailgun)

	opts := buildMailgunDomainOpts(d)

	log.Printf("[DEBUG] Domain create configuration: %#v", opts)

	err := client.CreateDomain(opts.Name, opts.SMTPPassword, opts.SpamAction, opts.Wildcard)
	if err != nil {
		return err
	}

	d.SetId(opts.Name)

	log.Printf("[INFO] Domain ID: %s", d.Id())

	// Retrieve and update state of domain
	err = resourceMailgunDomainRetrieve(d.Id(), client, d)

	if err != nil {
		return err
	}

	return nil
}

func resourceMailgunDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(mailgun.Mailgun)

	log.Printf("[INFO] Deleting Domain: %s", d.Id())

	// Destroy the domain
	err := client.DeleteDomain(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting domain: %s", err)
	}

	// Give the destroy a chance to take effect
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, _, _, err = client.GetSingleDomain(d.Id())
		if err == nil {
			log.Printf("[INFO] Retrying until domain disappears...")
			return resource.RetryableError(
				fmt.Errorf("Domain seems to still exist; will check again."))
		}
		log.Printf("[INFO] Got error looking for domain, seems gone: %s", err)
		return nil
	})
}

func resourceMailgunDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(mailgun.Mailgun)

	err := resourceMailgunDomainRetrieve(d.Id(), client, d)

	if err != nil {
		return err
	}

	return nil
}

func resourceMailgunDomainRetrieve(id string, client mailgun.Mailgun, d *schema.ResourceData) error {
	domain, _, _, err := client.GetSingleDomain(id)

	if err != nil {
		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	d.Set("name", domain.Name)
	d.Set("smtp_password", domain.SMTPPassword)
	d.Set("smtp_login", domain.SMTPLogin)
	d.Set("wildcard", domain.Wildcard)
	d.Set("spam_action", domain.SpamAction)

	return nil
}

func buildMailgunDomainOpts(d *schema.ResourceData) *mailgunDomainOpts {
	opts := &mailgunDomainOpts{
		Name:         d.Get("name").(string),
		SMTPPassword: d.Get("smtp_password").(string),
		SpamAction:   d.Get("spam_action").(string),
		Wildcard:     d.Get("wildcard").(bool),
	}

	return opts
}

type mailgunDomainOpts struct {
	Name         string
	SMTPPassword string
	SpamAction   string
	Wildcard     bool
}
