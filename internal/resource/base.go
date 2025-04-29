package resource

import (
	"encoding/json"
	"fmt"

	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BaseResource represents the common functionality for all Aidbox resources
type BaseResource struct {
	ResourceType string
	Schema       map[string]*schema.Schema
	CreateFunc   func(d *schema.ResourceData, m interface{}) error
	ReadFunc     func(d *schema.ResourceData, m interface{}) error
	UpdateFunc   func(d *schema.ResourceData, m interface{}) error
	DeleteFunc   func(d *schema.ResourceData, m interface{}) error
}

// NewBaseResource creates a new base resource with common schema fields
func NewBaseResource(resourceType string) *BaseResource {
	return &BaseResource{
		ResourceType: resourceType,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			// Extensions field for additional properties
			"extensions": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Additional fields that are not explicitly defined in the schema",
			},
		},
		CreateFunc: ResourceBaseCreate,
		ReadFunc:   ResourceBaseRead,
		UpdateFunc: ResourceBaseUpdate,
		DeleteFunc: ResourceBaseDelete,
	}
}

// AddSchema adds a new schema field to the base resource
func (b *BaseResource) AddSchema(name string, schema *schema.Schema) {
	b.Schema[name] = schema
}

// SetCreateFunc sets a custom create function
func (b *BaseResource) SetCreateFunc(f func(d *schema.ResourceData, m interface{}) error) {
	b.CreateFunc = f
}

// SetReadFunc sets a custom read function
func (b *BaseResource) SetReadFunc(f func(d *schema.ResourceData, m interface{}) error) {
	b.ReadFunc = f
}

// SetUpdateFunc sets a custom update function
func (b *BaseResource) SetUpdateFunc(f func(d *schema.ResourceData, m interface{}) error) {
	b.UpdateFunc = f
}

// SetDeleteFunc sets a custom delete function
func (b *BaseResource) SetDeleteFunc(f func(d *schema.ResourceData, m interface{}) error) {
	b.DeleteFunc = f
}

// ToResource converts the base resource to a schema.Resource
func (b *BaseResource) ToResource() *schema.Resource {
	return &schema.Resource{
		Create: b.CreateFunc,
		Read:   b.ReadFunc,
		Update: b.UpdateFunc,
		Delete: b.DeleteFunc,
		Schema: b.Schema,
	}
}

// ResourceBaseCreate handles the creation of a new Aidbox resource
func ResourceBaseCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	// Get the resource ID
	resourceID := d.Get("resource_id").(string)
	if resourceID == "" {
		resourceID = fmt.Sprintf("tf-%s", d.Get("id").(string))
	}

	// Create a map for the resource
	resourceMap := make(map[string]interface{})
	resourceMap["resourceType"] = d.Get("resource_type").(string)
	resourceMap["id"] = resourceID

	// Add extensions if provided
	if extensions, ok := d.GetOk("extensions"); ok {
		extensionsMap := extensions.(map[string]interface{})
		for k, v := range extensionsMap {
			resourceMap[k] = v
		}
	}

	// Convert to JSON
	resourceJSON, err := json.Marshal(resourceMap)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	if err := client.CreateResource(d.Get("resource_type").(string), resourceID, string(resourceJSON)); err != nil {
		return err
	}

	d.SetId(resourceID)
	return ResourceBaseRead(d, m)
}

// ResourceBaseRead handles reading an existing Aidbox resource
func ResourceBaseRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Id()
	resourceType := d.Get("resource_type").(string)
	resource, err := client.GetResource(resourceType, resourceID)
	if err != nil {
		return err
	}

	if resource == "" {
		d.SetId("")
		return nil
	}

	// Parse the resource JSON
	var resourceMap map[string]interface{}
	if err := json.Unmarshal([]byte(resource), &resourceMap); err != nil {
		return fmt.Errorf("failed to parse resource JSON: %w", err)
	}

	// Set the meta field if it exists
	if meta, ok := resourceMap["meta"].(map[string]interface{}); ok {
		metaList := []map[string]interface{}{make(map[string]interface{})}
		metaMap := metaList[0]

		// Handle simple string fields
		if v, ok := meta["version_id"].(string); ok {
			metaMap["version_id"] = v
		}
		if v, ok := meta["last_updated"].(string); ok {
			metaMap["last_updated"] = v
		}
		if v, ok := meta["created_at"].(string); ok {
			metaMap["created_at"] = v
		}

		d.Set("meta", metaList)
	}

	// Set resource_type
	d.Set("resource_type", resourceType)

	// Set extensions for fields that aren't explicitly defined in the schema
	extensions := make(map[string]string)
	for k, v := range resourceMap {
		// Skip fields that are explicitly defined in the schema
		if _, ok := d.GetOk(k); ok {
			continue
		}
		// Skip meta and id fields
		if k == "meta" || k == "id" || k == "resourceType" {
			continue
		}
		// Convert value to string
		if str, ok := v.(string); ok {
			extensions[k] = str
		} else if str, ok := v.(float64); ok {
			extensions[k] = fmt.Sprintf("%v", str)
		} else if str, ok := v.(bool); ok {
			extensions[k] = fmt.Sprintf("%v", str)
		} else if str, ok := v.([]interface{}); ok {
			// Handle arrays
			jsonBytes, err := json.Marshal(str)
			if err == nil {
				extensions[k] = string(jsonBytes)
			}
		} else if str, ok := v.(map[string]interface{}); ok {
			// Handle nested objects
			jsonBytes, err := json.Marshal(str)
			if err == nil {
				extensions[k] = string(jsonBytes)
			}
		}
	}
	d.Set("extensions", extensions)

	return nil
}

// ResourceBaseUpdate handles updating an existing Aidbox resource
func ResourceBaseUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Id()
	resourceType := d.Get("resource_type").(string)

	// Create a map for the resource
	resourceMap := make(map[string]interface{})
	resourceMap["resourceType"] = resourceType
	resourceMap["id"] = resourceID

	// Add extensions if provided
	if extensions, ok := d.GetOk("extensions"); ok {
		extensionsMap := extensions.(map[string]interface{})
		for k, v := range extensionsMap {
			resourceMap[k] = v
		}
	}

	// Convert to JSON
	resourceJSON, err := json.Marshal(resourceMap)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	if err := client.UpdateResource(resourceType, resourceID, string(resourceJSON)); err != nil {
		return err
	}

	return ResourceBaseRead(d, m)
}

// ResourceBaseDelete handles deleting an existing Aidbox resource
func ResourceBaseDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Id()
	resourceType := d.Get("resource_type").(string)
	return client.DeleteResource(resourceType, resourceID)
}
