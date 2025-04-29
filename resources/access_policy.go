package resources

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/flawless/terraform-provider-aidbox/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAidboxAccessPolicy() *schema.Resource {
	base := resource.NewBaseResource("AccessPolicy")

	// Add access policy-specific schema fields
	base.AddSchema("engine", &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"json-schema",
			"allow",
			"sql",
			"complex",
			"matcho",
			"clj",
			"matcho-rpc",
			"allow-rpc",
			"signed-rpc",
			"smart-on-fhir",
		}, false),
		Description: "The engine for the access policy (json-schema, allow, sql, complex, matcho, clj, matcho-rpc, allow-rpc, signed-rpc, smart-on-fhir)",
	})

	// Schema field for matcho engine
	base.AddSchema("matcho", &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "Matcho engine configuration",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	})

	// Schema field for SQL engine
	base.AddSchema("sql", &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "SQL engine configuration",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	})

	// Schema field for JSON Schema engine
	base.AddSchema("schema", &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "JSON Schema engine configuration",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	})

	// Schema field for Complex engine
	base.AddSchema("and", &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Complex engine AND rules",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	})

	base.AddSchema("or", &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Complex engine OR rules",
		Elem: &schema.Schema{
			Type: schema.TypeString,
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

		// Handle engine-specific configurations
		engine := d.Get("engine").(string)
		switch engine {
		case "matcho":
			if v, ok := d.GetOk("matcho"); ok {
				accessPolicyMap["matcho"] = v
			}
		case "sql":
			if v, ok := d.GetOk("sql"); ok {
				accessPolicyMap["sql"] = v
			}
		case "json-schema":
			if v, ok := d.GetOk("schema"); ok {
				accessPolicyMap["schema"] = v
			}
		case "complex":
			if v, ok := d.GetOk("and"); ok {
				accessPolicyMap["and"] = v
			}
			if v, ok := d.GetOk("or"); ok {
				accessPolicyMap["or"] = v
			}
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

		// Handle engine-specific configurations
		engine := d.Get("engine").(string)
		switch engine {
		case "matcho":
			if v, ok := d.GetOk("matcho"); ok {
				accessPolicyMap["matcho"] = v
			}
		case "sql":
			if v, ok := d.GetOk("sql"); ok {
				accessPolicyMap["sql"] = v
			}
		case "json-schema":
			if v, ok := d.GetOk("schema"); ok {
				accessPolicyMap["schema"] = v
			}
		case "complex":
			if v, ok := d.GetOk("and"); ok {
				accessPolicyMap["and"] = v
			}
			if v, ok := d.GetOk("or"); ok {
				accessPolicyMap["or"] = v
			}
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
