package linode

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeInstanceTypeRead,

		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deprecated": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_public": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vendor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLinodeImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*linodego.Client)

	images, err := client.ListImages(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error listing images: %s", err)
	}

	reqImage := d.Get("id").(string)

	for _, r := range images {
		if r.ID == reqImage {
			d.SetId(r.ID)
			d.Set("label", r.Label)
			d.Set("description", r.Description)
			d.Set("created", r.Created)
			d.Set("created_by", r.CreatedBy)
			d.Set("deprecated", r.Deprecated)
			d.Set("is_public", r.IsPublic)
			d.Set("size", r.Size)
			d.Set("type", r.Type)
			d.Set("vendor", r.Vendor)
			return nil
		}
	}

	d.SetId("")

	return fmt.Errorf("Image %s was not found", reqImage)
}
