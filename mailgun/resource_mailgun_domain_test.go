package mailgun

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	mailgun "github.com/mailgun/mailgun-go/v3"
)

var _testDomainName = "terraformv3.example.com"

func TestAccMailgunDomain_Basic(t *testing.T) {
	var resp mailgun.DomainResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunDomainConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunDomainExists("mailgun_domain.foobar", &resp),
					testAccCheckMailgunDomainAttributes(&resp),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "name", _testDomainName),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "spam_action", "disabled"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "wildcard", "true"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.0.priority", "10"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.0.name", _testDomainName),
				),
			},
		},
	})
}

func TestAccMailgunDomain_Import(t *testing.T) {
	resourceName := "mailgun_domain.foobar"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunDomainConfig(),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMailgunDomainDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mailgun_domain" {
			continue
		}

		client, errc := testAccProvider.Meta().(*Config).GetClient(rs.Primary.Attributes["region"])
		if errc != nil {
			return errc
		}

		resp, err := client.GetDomain(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists: %#v", resp)
		}
	}

	return nil
}

func testAccCheckMailgunDomainAttributes(DomainResp *mailgun.DomainResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if DomainResp.Domain.Name != _testDomainName {
			return fmt.Errorf("Bad name: %s", DomainResp.Domain.Name)
		}

		if DomainResp.Domain.SpamAction != "disabled" {
			return fmt.Errorf("Bad spam_action: %s", DomainResp.Domain.SpamAction)
		}

		if DomainResp.Domain.Wildcard != true {
			return fmt.Errorf("Bad wildcard: %t", DomainResp.Domain.Wildcard)
		}

		if DomainResp.ReceivingDNSRecords[0].Priority == "" {
			return fmt.Errorf("Bad receiving_records: %s", DomainResp.ReceivingDNSRecords)
		}

		if DomainResp.SendingDNSRecords[0].Name == "" {
			return fmt.Errorf("Bad sending_records: %s", DomainResp.SendingDNSRecords)
		}

		return nil
	}
}

func testAccCheckMailgunDomainExists(n string, DomainResp *mailgun.DomainResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain ID is set")
		}

		client, errc := testAccProvider.Meta().(*Config).GetClient(rs.Primary.Attributes["region"])
		if errc != nil {
			return errc
		}

		resp, err := client.GetDomain(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if resp.Domain.Name != rs.Primary.ID {
			return fmt.Errorf("Domain not found")
		}

		*DomainResp = resp

		return nil
	}
}

func testAccCheckMailgunDomainConfig() string {
	return `resource "mailgun_domain" "foobar" {
    name = "` + _testDomainName + `"
	spam_action = "disabled"
	region = "us"
    wildcard = true
}`
}
