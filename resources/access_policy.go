package resources

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxAccessPolicy() *schema.Resource {
	base := resource.NewBaseResource("AccessPolicy")

	// Add access policy-specific schema fields
	base.AddSchema("engine", &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The engine for the access policy",
	})

	base.AddSchema("rules", &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The rules for the access policy",
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	})

	// Override the create function to handle the access policy-specific fields
	base.SetCreateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Set the resource type
		d.Set("resource_type", "AccessPolicy")

		// Call the base create function to handle common fields
		if err := resource.ResourceBaseCreate(d, m); err != nil {
			return err
		}

		// Get the resource ID
		resourceID := d.Id()

		// Create a map for the access policy resource
		accessPolicyMap := map[string]interface{}{
			"resourceType": "AccessPolicy",
			"id":           resourceID,
			"engine":       d.Get("engine").(string),
		}

		// Handle rules - directly use the rules as provided
		if rules := d.Get("rules").([]interface{}); len(rules) > 0 {
			accessPolicyMap["rules"] = rules
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				accessPolicyMap[k] = v
			}
		}

		// Convert to JSON
		accessPolicyJSON, err := json.Marshal(accessPolicyMap)
		if err != nil {
			return fmt.Errorf("failed to marshal access policy: %w", err)
		}

		// Update the resource with the access policy-specific fields
		client := m.(*client.Client)
		if err := client.UpdateResource("AccessPolicy", resourceID, string(accessPolicyJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	// Override the update function to handle the access policy-specific fields
	base.SetUpdateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Get the resource ID
		resourceID := d.Id()

		// Create a map for the access policy resource
		accessPolicyMap := map[string]interface{}{
			"resourceType": "AccessPolicy",
			"id":           resourceID,
			"engine":       d.Get("engine").(string),
		}

		// Handle rules - directly use the rules as provided
		if rules := d.Get("rules").([]interface{}); len(rules) > 0 {
			accessPolicyMap["rules"] = rules
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				accessPolicyMap[k] = v
			}
		}

		// Convert to JSON
		accessPolicyJSON, err := json.Marshal(accessPolicyMap)
		if err != nil {
			return fmt.Errorf("failed to marshal access policy: %w", err)
		}

		// Update the resource with the access policy-specific fields
		client := m.(*client.Client)
		if err := client.UpdateResource("AccessPolicy", resourceID, string(accessPolicyJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	return base.ToResource()
}
