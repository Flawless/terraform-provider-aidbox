package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAidboxUser(t *testing.T) {
	resourceName := acctest.RandString(8)
	givenName := acctest.RandString(8)
	familyName := "Test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("AIDBOX_URL"); v == "" {
				t.Fatal("AIDBOX_URL must be set for acceptance tests")
			}
			if v := os.Getenv("AIDBOX_CLIENT_ID"); v == "" {
				t.Fatal("AIDBOX_CLIENT_ID must be set for acceptance tests")
			}
			if v := os.Getenv("AIDBOX_CLIENT_SECRET"); v == "" {
				t.Fatal("AIDBOX_CLIENT_SECRET must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"aidbox": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxUserConfig(resourceName, givenName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAidboxUserExists("aidbox_user.test"),
					testAccCheckAidboxUserAttributes("aidbox_user.test", givenName, familyName),
				),
			},
		},
	})
}

func testAccCheckAidboxUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No User ID is set")
		}

		return nil
	}
}

func testAccCheckAidboxUserAttributes(n string, expectedGivenName string, expectedFamilyName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		givenName := rs.Primary.Attributes["name.0.given_name"]
		if givenName != expectedGivenName {
			return fmt.Errorf("Expected given_name %s, got %s", expectedGivenName, givenName)
		}

		familyName := rs.Primary.Attributes["name.0.family_name"]
		if familyName != expectedFamilyName {
			return fmt.Errorf("Expected family_name %s, got %s", expectedFamilyName, familyName)
		}

		return nil
	}
}

func testAccAidboxUserConfig(resourceID string, givenName string, familyName string) string {
	return fmt.Sprintf(`
provider "aidbox" {
  url           = "%s"
  client_id     = "%s"
  client_secret = "%s"
}

resource "aidbox_user" "test" {
  resource_id = "%s"
  name {
    given_name  = "%s"
    family_name = "%s"
  }
  password    = "testpassword123"
}
`,
		os.Getenv("AIDBOX_URL"),
		os.Getenv("AIDBOX_CLIENT_ID"),
		os.Getenv("AIDBOX_CLIENT_SECRET"),
		resourceID,
		givenName,
		familyName,
	)
}
