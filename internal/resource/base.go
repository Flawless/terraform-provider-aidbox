package resource

import (
	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// BaseResource represents the common functionality for all Aidbox resources
type BaseResource struct {
	ResourceType string
	Schema       map[string]*schema.Schema
}

// NewBaseResource creates a new base resource with common schema fields
func NewBaseResource(resourceType string) *BaseResource {
	return &BaseResource{
		ResourceType: resourceType,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"meta": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// AddSchema adds a new schema field to the base resource
func (b *BaseResource) AddSchema(name string, schema *schema.Schema) {
	b.Schema[name] = schema
}

// Create handles the creation of a new Aidbox resource
func (b *BaseResource) Create(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Get("id").(string)
	resource := d.Get("resource").(string)

	if err := client.CreateResource(b.ResourceType, resourceID, resource); err != nil {
		return err
	}

	d.SetId(resourceID)
	return b.Read(d, m)
}

// Read handles reading an existing Aidbox resource
func (b *BaseResource) Read(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Get("id").(string)
	resource, err := client.GetResource(b.ResourceType, resourceID)
	if err != nil {
		return err
	}

	if resource == "" {
		d.SetId("")
		return nil
	}

	d.Set("resource", resource)
	d.Set("resource_type", b.ResourceType)
	return nil
}

// Update handles updating an existing Aidbox resource
func (b *BaseResource) Update(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Get("id").(string)
	resource := d.Get("resource").(string)

	if err := client.UpdateResource(b.ResourceType, resourceID, resource); err != nil {
		return err
	}

	return b.Read(d, m)
}

// Delete handles deleting an existing Aidbox resource
func (b *BaseResource) Delete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client.Client)

	resourceID := d.Get("id").(string)
	return client.DeleteResource(b.ResourceType, resourceID)
}
