package mailgun

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMailgunCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceMailgunCredentialCreate,
		Read:   resourceMailgunCredentialRead,
		Delete: resourceMailgunCredentialDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMailgunCredentialImporter,
		},

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"login": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
				Optional: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceMailgunCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*MailgunProvider).Domain(d.Get("domain").(string))

	opts := mailgun.Credential{}

	log.Printf("[DEBUG] Credential create configuration: %#v", opts)

	err := client.CreateCredential(d.Get("login").(string), d.Get("password").(string))

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("domain").(string), d.Get("login").(string)))

	log.Printf("[INFO] Credential ID: %s", d.Id())

	// Retrieve and update state of Credential
	_, err = resourceMailginCredentialRetrieve(d.Id(), client, d)

	if err != nil {
		return err
	}

	return nil
}

func resourceMailgunCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*MailgunProvider).Domain(d.Get("domain").(string))

	log.Printf("[INFO] Deleting Credential: %s", d.Id())

	// Destroy the Credential
	err := client.DeleteCredential(d.Get("login").(string))
	if err != nil {
		return fmt.Errorf("Error deleting Credential: %s", err)
	}
	return nil
}

func resourceMailgunCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*MailgunProvider).Domain(d.Get("domain").(string))

	_, err := resourceMailginCredentialRetrieve(d.Id(), client, d)

	if err != nil {
		return err
	}

	return nil
}

func resourceMailginCredentialRetrieve(id string, client mailgun.Mailgun, d *schema.ResourceData) (*mailgun.Credential, error) {
	_, credentials, err := client.GetCredentials(mailgun.DefaultLimit, mailgun.DefaultSkip)

	if err != nil {
		return nil, fmt.Errorf("Error retrieving Credential: %s", err)
	}
	var resp mailgun.Credential

	for _, cred := range credentials {
		if cred.Login == d.Get("login").(string) {
			resp = cred
		}

	}

	d.Set("login", resp.Login)

	return &resp, nil
}

func resourceMailgunCredentialImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid mailgun credential specifier. Expecting {domain}/{login}")
	}
	d.Set("domain", parts[0])
	d.Set("login", parts[1])
	return []*schema.ResourceData{d}, nil
}
