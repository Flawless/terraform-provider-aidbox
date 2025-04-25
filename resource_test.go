package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAidboxResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAidboxResourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"aidbox_resource.test", "resource_type", "Organization",
					),
					resource.TestCheckResourceAttr(
						"aidbox_resource.test", "resource_id", "test-org",
					),
				),
			},
		},
	})
}

const testAccAidboxResourceConfig_basic = `
resource "aidbox_resource" "test" {
  resource_type = "Organization"
  resource_id   = "test-org"
  resource = jsonencode({
    name = "Test Organization"
  })
}
`
