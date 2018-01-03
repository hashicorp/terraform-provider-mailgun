package mailgun

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pearkes/mailgun"
)

func dataSourceMailgunDomain() *schema.Resource {

	mailgunSchema := resourceMailgunSchema()
	mailgunSchema["smtp_password"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}

	return &schema.Resource{
		Read:   dataSourceMailgunDomainRead,
		Schema: mailgunSchema,
	}
}

func dataSourceMailgunDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mailgun.Client)
	name := d.Get("name").(string)

	_, err := resourceMailginDomainRetrieve(name, client, d)

	if err != nil {
		return err
	}

	d.SetId(name)
	return nil
}
