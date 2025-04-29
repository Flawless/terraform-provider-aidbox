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

func TestAccAidboxRole(t *testing.T) {
	resourceName := acctest.RandString(8)
	roleName := "test-role"
	userName := "test-user"

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
				Config: testAccAidboxRoleConfig(resourceName, roleName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAidboxRoleExists("aidbox_role.test"),
					testAccCheckAidboxRoleAttributes("aidbox_role.test", roleName, userName),
				),
			},
		},
	})
}

func TestAccAidboxAccessPolicy(t *testing.T) {
	resourceName := acctest.RandString(8)
	engine := "matcho"

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
				Config: testAccAidboxAccessPolicyConfig(resourceName, engine),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAidboxAccessPolicyExists("aidbox_access_policy.test"),
					testAccCheckAidboxAccessPolicyAttributes("aidbox_access_policy.test", engine),
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

func testAccCheckAidboxRoleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Role ID is set")
		}

		return nil
	}
}

func testAccCheckAidboxRoleAttributes(n string, expectedName string, expectedUser string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		name := rs.Primary.Attributes["name"]
		if name != expectedName {
			return fmt.Errorf("Expected name %s, got %s", expectedName, name)
		}

		userID := rs.Primary.Attributes["user.0.id"]
		if userID != expectedUser {
			return fmt.Errorf("Expected user ID %s, got %s", expectedUser, userID)
		}

		return nil
	}
}

func testAccCheckAidboxAccessPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AccessPolicy ID is set")
		}

		return nil
	}
}

func testAccCheckAidboxAccessPolicyAttributes(n string, expectedEngine string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		engine := rs.Primary.Attributes["engine"]
		if engine != expectedEngine {
			return fmt.Errorf("Expected engine %s, got %s", expectedEngine, engine)
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

func testAccAidboxRoleConfig(resourceID string, name string, user string) string {
	return fmt.Sprintf(`
provider "aidbox" {
  url           = "%s"
  client_id     = "%s"
  client_secret = "%s"
}

resource "aidbox_user" "test_user" {
  resource_id = "%s"
  name {
    given_name  = "Test"
    family_name = "User"
  }
  password = "testpassword123"
}

resource "aidbox_role" "test" {
  resource_id = "%s"
  name        = "%s"
  user {
    id = aidbox_user.test_user.id
  }
  extensions = {
    "description" = "Test role created by Terraform"
  }
}
`,
		os.Getenv("AIDBOX_URL"),
		os.Getenv("AIDBOX_CLIENT_ID"),
		os.Getenv("AIDBOX_CLIENT_SECRET"),
		user,
		resourceID,
		name,
	)
}

func testAccAidboxAccessPolicyConfig(resourceID string, engine string) string {
	return fmt.Sprintf(`
provider "aidbox" {
  url           = "%s"
  client_id     = "%s"
  client_secret = "%s"
}

resource "aidbox_access_policy" "test" {
  resource_id = "%s"
  engine      = "%s"
  matcho = {
    "request-method" = "get"
    "uri"           = "/Patient"
  }
}
`,
		os.Getenv("AIDBOX_URL"),
		os.Getenv("AIDBOX_CLIENT_ID"),
		os.Getenv("AIDBOX_CLIENT_SECRET"),
		resourceID,
		engine,
	)
}
