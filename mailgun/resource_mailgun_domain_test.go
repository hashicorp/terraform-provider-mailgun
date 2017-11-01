package mailgun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/mailgun/mailgun-go.v1"
)

func TestAccMailgunDomain_Basic(t *testing.T) {
	var domain mailgun.Domain
	var receivingDNSRecords []mailgun.DNSRecord
	var sendingDNSRecords []mailgun.DNSRecord

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunDomainConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunDomainExists("mailgun_domain.foobar", &domain, &receivingDNSRecords, &sendingDNSRecords),
					testAccCheckMailgunDomainAttributes(&domain, &receivingDNSRecords, &sendingDNSRecords),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "name", "terraform.example.com"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "spam_action", "disabled"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "smtp_password", "foobar"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "wildcard", "true"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.0.priority", "10"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.0.name", "terraform.example.com"),
				),
			},
		},
	})
}

func testAccCheckMailgunDomainDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*MailgunProvider).Domain("")

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

func testAccCheckMailgunDomainAttributes(Domain *mailgun.Domain, ReceivingDNSRecords *[]mailgun.DNSRecord, SendingDNSRecords *[]mailgun.DNSRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if Domain.Name != "terraform.example.com" {
			return fmt.Errorf("Bad name: %s", Domain.Name)
		}

		if Domain.SpamAction != "disabled" {
			return fmt.Errorf("Bad spam_action: %s", Domain.SpamAction)
		}

		if Domain.Wildcard != true {
			return fmt.Errorf("Bad wildcard: %t", Domain.Wildcard)
		}

		if Domain.SMTPPassword != "foobar" {
			return fmt.Errorf("Bad smtp_password: %s", Domain.SMTPPassword)
		}

		if (*ReceivingDNSRecords)[0].Priority == "" {
			return fmt.Errorf("Bad receiving_records: %s", *ReceivingDNSRecords)
		}

		if (*SendingDNSRecords)[0].Name == "" {
			return fmt.Errorf("Bad sending_records: %s", *SendingDNSRecords)
		}

		return nil
	}
}

func testAccCheckMailgunDomainExists(n string, Domain *mailgun.Domain, ReceivingDNSRecords *[]mailgun.DNSRecord, SendingDNSRecords *[]mailgun.DNSRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client := testAccProvider.Meta().(*MailgunProvider).Domain(rs.Primary.ID)

		var err error
		*Domain, *ReceivingDNSRecords, *SendingDNSRecords, err = client.GetSingleDomain(rs.Primary.ID)

		if err != nil {
			return err
		}

		if Domain.Name != rs.Primary.ID {
			return fmt.Errorf("Domain not found")
		}

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
