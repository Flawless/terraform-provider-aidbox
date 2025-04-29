package resources

import (
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxRole() *schema.Resource {
	base := resource.NewBaseResource("Role")

	// Add role-specific schema fields
	base.AddSchema("type", &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	})

	return &schema.Resource{
		Create: base.Create,
		Read:   base.Read,
		Update: base.Update,
		Delete: base.Delete,
		Schema: base.Schema,
	}
}
