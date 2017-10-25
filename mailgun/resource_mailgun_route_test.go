package mailgun

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mailgun/mailgun-go.v1"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccMailgunRoute_Basic(t *testing.T) {
	var route mailgun.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMailgunRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckMailgunRouteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMailgunRouteExists("mailgun_route.foobar", &route),
					resource.TestCheckResourceAttr(
						"mailgun_route.foobar", "priority", "0"),
					resource.TestCheckResourceAttr(
						"mailgun_route.foobar", "description", "inbound"),
					resource.TestCheckResourceAttr(
						"mailgun_route.foobar", "expression", "match_recipient('.*@example.com')"),
					resource.TestCheckResourceAttr(
						"mailgun_route.foobar", "actions.0", "forward('http://example.com/api/v1/foos/')"),
					resource.TestCheckResourceAttr(
						"mailgun_route.foobar", "actions.1", "stop()"),
				),
			},
		},
	})
}

func testAccCheckMailgunRouteDestroy(s *terraform.State) error {
	client := *testAccProvider.Meta().(*mailgun.Mailgun)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mailgun_route" {
			continue
		}

		route, err := client.GetRouteByID(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Route still exists: %#v", route)
		}
	}

	return nil
}

func testAccCheckMailgunRouteExists(n string, Route *mailgun.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Route ID is set")
		}

		client := *testAccProvider.Meta().(*mailgun.Mailgun)

		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			var err error
			*Route, err = client.GetRouteByID(rs.Primary.ID)

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("Unable to find Route after retries: %s", err)
		}

		if Route.ID != rs.Primary.ID {
			return fmt.Errorf("Route not found")
		}

		return nil
	}
}

const testAccCheckMailgunRouteConfig_basic = `
resource "mailgun_route" "foobar" {
    priority = "0"
    description = "inbound"
    expression = "match_recipient('.*@example.com')"
    actions = [
        "forward('http://example.com/api/v1/foos/')",
        "stop()"
    ]
}
`
