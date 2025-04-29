package resources

import (
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxUser() *schema.Resource {
	base := resource.NewBaseResource("User")

	// Add user-specific schema fields
	base.AddSchema("password", &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Sensitive: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return true
		},
	})

	return &schema.Resource{
		Create: base.Create,
		Read:   base.Read,
		Update: base.Update,
		Delete: base.Delete,
		Schema: base.Schema,
	}
}
