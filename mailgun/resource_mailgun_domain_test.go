package mailgun

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/pearkes/mailgun"
)

func TestAccMailgunDomain_Basic(t *testing.T) {
	var resp mailgun.DomainResponse
	uuid, _ := uuid.GenerateUUID()
	domain := fmt.Sprintf("terraform.%s.com", uuid)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunDomainDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunDomainConfig_basic(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunDomainExists("mailgun_domain.foobar", &resp),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "name", domain),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "spam_action", "disabled"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "smtp_password", "foobar"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "wildcard", "true"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.0.priority", "10"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.0.value", "mxa.mailgun.org"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.1.priority", "10"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "receiving_records.1.value", "mxb.mailgun.org"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.1.name", fmt.Sprintf("email.%s", domain)),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.1.value", "mailgun.org"),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.2.name", domain),
					resource.TestCheckResourceAttr(
						"mailgun_domain.foobar", "sending_records.2.value", "v=spf1 include:mailgun.org ~all"),
					resource.TestMatchResourceAttr(
						"mailgun_domain.foobar", "sending_records.0.name", regexp.MustCompile(fmt.Sprintf("_domainkey.%s$", domain))),
					resource.TestMatchResourceAttr(
						"mailgun_domain.foobar", "sending_records.0.value", regexp.MustCompile("^k=rsa; p=")),
				),
			},
		},
	})
}

func testAccCheckMailgunDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*mailgun.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mailgun_domain" {
			continue
		}

		resp, err := client.RetrieveDomain(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists: %#v", resp)
		}
	}

	return nil
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

		client := testAccProvider.Meta().(*mailgun.Client)

		resp, err := client.RetrieveDomain(rs.Primary.ID)

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

func testAccCheckMailgunDomainConfig_basic(domain string) string {
	return fmt.Sprintf(`
resource "mailgun_domain" "foobar" {
    name = "%s"
    spam_action = "disabled"
    smtp_password = "foobar"
    wildcard = true
}`, domain)
}
