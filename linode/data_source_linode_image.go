package linode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeImageRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLinodeImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	reqImage := d.Get("id").(string)

	if reqImage == "" {
		return fmt.Errorf("Image id is required")
	}

	image, err := client.GetImage(context.Background(), reqImage)
	if err != nil {
		return fmt.Errorf("Error listing images: %s", err)
	}

	if image != nil {
		d.SetId(image.ID)
		d.Set("label", image.Label)
		d.Set("description", image.Description)
		d.Set("created", image.Created)
		d.Set("created_by", image.CreatedBy)
		d.Set("deprecated", image.Deprecated)
		d.Set("is_public", image.IsPublic)
		d.Set("size", image.Size)
		d.Set("type", image.Type)
		d.Set("vendor", image.Vendor)
		return nil
	}

	d.SetId("")

	return fmt.Errorf("Image %s was not found", reqImage)
}
