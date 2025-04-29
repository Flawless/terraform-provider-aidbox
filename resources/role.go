package resources

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxRole() *schema.Resource {
	base := resource.NewBaseResource("Role")

	// Add role-specific schema fields
	base.AddSchema("name", &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the role",
	})

	base.AddSchema("user", &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The ID of the user",
				},
				"resource_type": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "User",
					Description: "The resource type of the user reference",
				},
			},
		},
		Description: "The user reference for the role",
	})

	// Override the create function to handle the role-specific fields
	base.SetCreateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Set the resource type
		d.Set("resource_type", "Role")

		// Get or generate the resource ID
		resourceID := d.Get("resource_id").(string)
		if resourceID == "" {
			resourceID = fmt.Sprintf("tf-%s", d.Get("id").(string))
		}

		// Create a map for the role resource
		roleMap := map[string]interface{}{
			"resourceType": "Role",
			"id":           resourceID,
			"name":         d.Get("name").(string),
		}

		// Handle user reference
		if v, ok := d.GetOk("user"); ok {
			userList := v.([]interface{})
			if len(userList) > 0 {
				userMap := userList[0].(map[string]interface{})
				roleMap["user"] = map[string]interface{}{
					"id":           userMap["id"].(string),
					"resourceType": userMap["resource_type"].(string),
				}
			}
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

		// Create the resource
		client := m.(*client.Client)
		if err := client.CreateResource("Role", resourceID, string(roleJSON)); err != nil {
			return err
		}

		d.SetId(resourceID)
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
			"name":         d.Get("name").(string),
		}

		// Handle user reference
		if v, ok := d.GetOk("user"); ok {
			userList := v.([]interface{})
			if len(userList) > 0 {
				userMap := userList[0].(map[string]interface{})
				roleMap["user"] = map[string]interface{}{
					"id":           userMap["id"].(string),
					"resourceType": userMap["resource_type"].(string),
				}
			}
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

		// Update the resource
		client := m.(*client.Client)
		if err := client.UpdateResource("Role", resourceID, string(roleJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	return base.ToResource()
}
