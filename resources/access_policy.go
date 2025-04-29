package resources

import (
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxAccessPolicy() *schema.Resource {
	base := resource.NewBaseResource("AccessPolicy")

	// Add access policy-specific schema fields
	base.AddSchema("engine", &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	})

	base.AddSchema("rules", &schema.Schema{
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
