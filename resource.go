package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAidboxResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAidboxResourceCreate,
		Read:   resourceAidboxResourceRead,
		Update: resourceAidboxResourceUpdate,
		Delete: resourceAidboxResourceDelete,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAidboxResourceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	resourceID := d.Get("resource_id").(string)
	resource := d.Get("resource").(string)

	err := client.CreateResource(resourceType, resourceID, resource)
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", resourceType, resourceID))
	return resourceAidboxResourceRead(d, meta)
}

// filterServerFields removes Aidbox server-generated fields from resource JSON
func filterServerFields(resource string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(resource), &obj); err != nil {
		return "", err
	}
	delete(obj, "meta")
	if pw, ok := obj["password"].(string); ok && strings.HasPrefix(pw, "$s0$") {
		delete(obj, "password")
	}
	filtered, err := json.Marshal(obj)
	return string(filtered), err
}

func resourceAidboxResourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	resourceID := d.Get("resource_id").(string)

	resource, err := client.GetResource(resourceType, resourceID)
	if err != nil {
		return fmt.Errorf("error reading resource: %w", err)
	}

	if resource == "" {
		d.SetId("")
		return nil
	}

	filtered, err := filterServerFields(resource)
	if err != nil {
		return fmt.Errorf("error filtering server fields: %w", err)
	}
	d.Set("resource", filtered)
	return nil
}

func resourceAidboxResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	resourceID := d.Get("resource_id").(string)
	resource := d.Get("resource").(string)

	err := client.UpdateResource(resourceType, resourceID, resource)
	if err != nil {
		return fmt.Errorf("error updating resource: %w", err)
	}

	return resourceAidboxResourceRead(d, meta)
}

func resourceAidboxResourceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	resourceID := d.Get("resource_id").(string)

	err := client.DeleteResource(resourceType, resourceID)
	if err != nil {
		return fmt.Errorf("error deleting resource: %w", err)
	}

	return nil
}
