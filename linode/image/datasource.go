package image

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	reqImage := d.Get("id").(string)

	if reqImage == "" {
		return diag.Errorf("Image id is required")
	}

	image, err := client.GetImage(ctx, reqImage)
	if err != nil {
		return diag.Errorf("Error listing images: %s", err)
	}

	if image != nil {
		d.SetId(image.ID)
		d.Set("capabilities", image.Capabilities)
		d.Set("label", image.Label)
		d.Set("description", image.Description)
		if image.Created != nil {
			d.Set("created", image.Created.Format(time.RFC3339))
		}
		if image.Expiry != nil {
			d.Set("expiry", image.Expiry.Format(time.RFC3339))
		}
		d.Set("created_by", image.CreatedBy)
		d.Set("deprecated", image.Deprecated)
		d.Set("is_public", image.IsPublic)
		d.Set("size", image.Size)
		d.Set("status", image.Status)
		d.Set("type", image.Type)
		d.Set("vendor", image.Vendor)
		return nil
	}

	d.SetId("")

	return diag.Errorf("Image %s was not found", reqImage)
}
