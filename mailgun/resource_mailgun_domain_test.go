package mailgun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mailgun/mailgun-go"
)

func TestAccMailgunDomain_Basic(t *testing.T) {
	var resp mailgun.Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunDomainConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunDomainExists("mailgun_domain.foobar", &resp),
					testAccCheckMailgunDomainAttributes(&resp),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "name", "terraform.example.com"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "spam_action", "disabled"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "smtp_password", "foobar"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "wildcard", "true"),
				),
			},
		},
	})
}

func testAccCheckMailgunDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(mailgun.Mailgun)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mailgun_domain" {
			continue
		}

		domain, _, _, err := client.GetSingleDomain(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists: %#v", domain)
		}
	}

	return nil
}

func testAccCheckMailgunDomainAttributes(DomainResp *mailgun.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if DomainResp.Name != "terraform.example.com" {
			return fmt.Errorf("Bad name: %s", DomainResp.Name)
		}

		if DomainResp.SpamAction != "disabled" {
			return fmt.Errorf("Bad spam_action: %s", DomainResp.SpamAction)
		}

		if DomainResp.Wildcard != true {
			return fmt.Errorf("Bad wildcard: %t", DomainResp.Wildcard)
		}

		if DomainResp.SMTPPassword != "foobar" {
			return fmt.Errorf("Bad smtp_password: %s", DomainResp.SMTPPassword)
		}

		return nil
	}
}

func testAccCheckMailgunDomainExists(n string, DomainResp *mailgun.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client := testAccProvider.Meta().(mailgun.Mailgun)

		domain, _, _, err := client.GetSingleDomain(rs.Primary.ID)

		if err != nil {
			return err
		}

		if domain.Name != rs.Primary.ID {
			return fmt.Errorf("Domain not found")
		}

		*DomainResp = domain

		return nil
	}
}

const testAccCheckMailgunDomainConfig_basic = `
resource "mailgun_domain" "foobar" {
    name = "terraform.example.com"
    spam_action = "disabled"
    smtp_password = "foobar"
    wildcard = true
}`
