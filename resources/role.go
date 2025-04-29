package resources

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxRole() *schema.Resource {
	base := resource.NewBaseResource("Role")

	// Add role-specific schema fields
	base.AddSchema("type", &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The type of the role",
	})

	// Override the create function to handle the role-specific fields
	base.SetCreateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Call the base create function to handle common fields
		if err := resource.ResourceBaseCreate(d, m); err != nil {
			return err
		}

		// Get the resource ID
		resourceID := d.Id()

		// Create a map for the role resource
		roleMap := map[string]interface{}{
			"resourceType": "Role",
			"id":           resourceID,
			"type":         d.Get("type").(string),
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				roleMap[k] = v
			}
		}

		// Convert to JSON
		roleJSON, err := json.Marshal(roleMap)
		if err != nil {
			return fmt.Errorf("failed to marshal role: %w", err)
		}

		// Update the resource with the role-specific fields
		client := m.(interface {
			UpdateResource(resourceType, id string, resourceJSON string) error
		})
		if err := client.UpdateResource("Role", resourceID, string(roleJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	// Override the update function to handle the role-specific fields
	base.SetUpdateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Get the resource ID
		resourceID := d.Id()

		// Create a map for the role resource
		roleMap := map[string]interface{}{
			"resourceType": "Role",
			"id":           resourceID,
			"type":         d.Get("type").(string),
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				roleMap[k] = v
			}
		}

		// Convert to JSON
		roleJSON, err := json.Marshal(roleMap)
		if err != nil {
			return fmt.Errorf("failed to marshal role: %w", err)
		}

		// Update the resource with the role-specific fields
		client := m.(interface {
			UpdateResource(resourceType, id string, resourceJSON string) error
		})
		if err := client.UpdateResource("Role", resourceID, string(roleJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	return base.ToResource()
}
