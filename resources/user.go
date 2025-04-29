package resources

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAidboxUser() *schema.Resource {
	base := resource.NewBaseResource("User")

	// Add user-specific schema fields
	base.AddSchema("name", &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"given_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The given name of the user",
				},
				"family_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The family name of the user",
				},
			},
		},
	})

	base.AddSchema("password", &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Sensitive:   true,
		Description: "The password for the user",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// Always suppress password differences since Aidbox hashes them
			return true
		},
	})

	// Override the create function to handle the user-specific fields
	base.SetCreateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Set the resource type
		d.Set("resource_type", "User")

		// Call the base create function to handle common fields
		if err := resource.ResourceBaseCreate(d, m); err != nil {
			return err
		}

		// Get the resource ID
		resourceID := d.Id()

		// Get the name from the schema
		nameList := d.Get("name").([]interface{})
		if len(nameList) == 0 {
			return fmt.Errorf("name block is required")
		}
		nameMap := nameList[0].(map[string]interface{})

		// Create a map for the user resource
		userMap := map[string]interface{}{
			"resourceType": "User",
			"id":           resourceID,
			"name": map[string]interface{}{
				"givenName":  nameMap["given_name"],
				"familyName": nameMap["family_name"],
			},
		}

		// Add password if provided
		if password, ok := d.GetOk("password"); ok {
			userMap["password"] = password.(string)
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				userMap[k] = v
			}
		}

		// Convert to JSON
		userJSON, err := json.Marshal(userMap)
		if err != nil {
			return fmt.Errorf("failed to marshal user: %w", err)
		}

		// Update the resource with the user-specific fields
		client := m.(*client.Client)
		if err := client.UpdateResource("User", resourceID, string(userJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	// Override the update function to handle the user-specific fields
	base.SetUpdateFunc(func(d *schema.ResourceData, m interface{}) error {
		// Get the resource ID
		resourceID := d.Id()

		// Get the name from the schema
		nameList := d.Get("name").([]interface{})
		if len(nameList) == 0 {
			return fmt.Errorf("name block is required")
		}
		nameMap := nameList[0].(map[string]interface{})

		// Create a map for the user resource
		userMap := map[string]interface{}{
			"resourceType": "User",
			"id":           resourceID,
			"name": map[string]interface{}{
				"givenName":  nameMap["given_name"],
				"familyName": nameMap["family_name"],
			},
		}

		// Add password if provided
		if password, ok := d.GetOk("password"); ok {
			userMap["password"] = password.(string)
		}

		// Add extensions if provided
		if extensions, ok := d.GetOk("extensions"); ok {
			extensionsMap := extensions.(map[string]interface{})
			for k, v := range extensionsMap {
				userMap[k] = v
			}
		}

		// Convert to JSON
		userJSON, err := json.Marshal(userMap)
		if err != nil {
			return fmt.Errorf("failed to marshal user: %w", err)
		}

		// Update the resource with the user-specific fields
		client := m.(*client.Client)
		if err := client.UpdateResource("User", resourceID, string(userJSON)); err != nil {
			return err
		}

		return resource.ResourceBaseRead(d, m)
	})

	return base.ToResource()
}
