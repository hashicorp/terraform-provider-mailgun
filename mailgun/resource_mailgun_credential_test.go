package mailgun

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gopkg.in/mailgun/mailgun-go.v1"
)

func TestAccMailgunCredential_import(t *testing.T) {
	t.Parallel()

	name := acctest.RandString(10)
	resourceName := "mailgun_credential." + name
	conf := testAccCheckMailgunCredentialConfig_basic(name, name)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunCredentialDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
func TestAccMailgunCredential_Basic(t *testing.T) {
	t.Parallel()

	var credential mailgun.Credential

	name := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunCredentialDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunCredentialConfig_basic(name, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunCredentialExists("mailgun_credential."+name, &credential),
					testAccCheckMailgunCredentialAttributes(name, &credential),
				),
			},
		},
	})
}

func testAccCheckMailgunCredentialDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mailgun_credential" {
			continue
		}

		a := strings.Split(rs.Primary.ID, "/")

		client := testAccProvider.Meta().(*MailgunProvider).Domain(a[0])
		err := client.DeleteCredential(a[1])

		if err == nil {
			return fmt.Errorf("Credential still exists: %#v", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMailgunCredentialAttributes(login string, Credential *mailgun.Credential) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if Credential.Login != login+"@terraform.example.com" {
			return fmt.Errorf("Bad login: %s", Credential.Login)
		}

		return nil
	}
}

func testAccCheckMailgunCredentialExists(n string, Credential *mailgun.Credential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Credential ID is set")
		}
		a := strings.Split(rs.Primary.ID, "/")

		client := testAccProvider.Meta().(*MailgunProvider).Domain(a[0])

		var err error
		var credentials []mailgun.Credential
		_, credentials, err = client.GetCredentials(mailgun.DefaultLimit, mailgun.DefaultSkip)
		if err != nil {
			return err
		}
		for _, cred := range credentials {
			if cred.Login == a[1] {
				*Credential = cred
				return nil
			}
		}

		return fmt.Errorf("Credential not found")
	}
}

func testAccCheckMailgunCredentialConfig_basic(name, login string) string {
	return fmt.Sprintf(`
resource "mailgun_domain" "foobar" {
    name = "terraform.example.com"
    spam_action = "disabled"
    smtp_password = "foobar"
    wildcard = true
}

resource "mailgun_credential" "%s" {
    domain = "${mailgun_domain.foobar.name}"
    login = "%s@terraform.example.com"
    password = "foobarpwd"
}`, name, login)
}
