package main

import (
	"fmt"

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
			"id": {
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
	id := d.Get("id").(string)
	resource := d.Get("resource").(string)

	err := client.CreateResource(resourceType, id, resource)
	if err != nil {
		return fmt.Errorf("error creating resource: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", resourceType, id))
	return resourceAidboxResourceRead(d, meta)
}

func resourceAidboxResourceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	id := d.Get("id").(string)

	resource, err := client.GetResource(resourceType, id)
	if err != nil {
		return fmt.Errorf("error reading resource: %w", err)
	}

	if resource == "" {
		d.SetId("")
		return nil
	}

	d.Set("resource", resource)
	return nil
}

func resourceAidboxResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	id := d.Get("id").(string)
	resource := d.Get("resource").(string)

	err := client.UpdateResource(resourceType, id, resource)
	if err != nil {
		return fmt.Errorf("error updating resource: %w", err)
	}

	return resourceAidboxResourceRead(d, meta)
}

func resourceAidboxResourceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client := NewClient(config)

	resourceType := d.Get("resource_type").(string)
	id := d.Get("id").(string)

	err := client.DeleteResource(resourceType, id)
	if err != nil {
		return fmt.Errorf("error deleting resource: %w", err)
	}

	return nil
}
