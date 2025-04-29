package main

import (
	"github.com/flawless/terraform-provider-aidbox/internal/client"
	"github.com/flawless/terraform-provider-aidbox/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AIDBOX_URL", nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AIDBOX_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("AIDBOX_CLIENT_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"aidbox_user":          resources.ResourceAidboxUser(),
			"aidbox_role":          resources.ResourceAidboxRole(),
			"aidbox_access_policy": resources.ResourceAidboxAccessPolicy(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			config := &client.Config{
				URL:          d.Get("url").(string),
				ClientID:     d.Get("client_id").(string),
				ClientSecret: d.Get("client_secret").(string),
			}
			return client.NewClient(config), nil
		},
	}
}
