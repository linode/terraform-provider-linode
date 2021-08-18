package linode

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"

	"context"
	"strconv"
)

func dataSourceLinodeImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinodeImagesRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema([]string{"deprecated", "is_public", "label", "size", "vendor"}),
			"images": {
				Type:        schema.TypeList,
				Description: "The returned list of Images.",
				Computed:    true,
				Elem:        image.DataSource(),
			},
		},
	}
}

func dataSourceLinodeImagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	filter, err := constructFilterString(d, imageValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	images, err := client.ListImages(ctx, &linodego.ListOptions{
		Filter: filter,
	})

	if err != nil {
		return diag.Errorf("failed to list linode images: %s", err)
	}

	imagesFlattened := make([]interface{}, len(images))
	for i, image := range images {
		imagesFlattened[i] = flattenLinodeImage(&image)
	}

	d.SetId(filter)
	d.Set("images", imagesFlattened)

	return nil
}

func flattenLinodeImage(image *linodego.Image) map[string]interface{} {
	result := make(map[string]interface{})

	result["id"] = image.ID
	result["label"] = image.Label
	result["description"] = image.Description
	result["created_by"] = image.CreatedBy
	result["deprecated"] = image.Deprecated
	result["is_public"] = image.IsPublic
	result["size"] = image.Size
	result["type"] = image.Type
	result["vendor"] = image.Vendor
	result["status"] = image.Status

	if image.Created != nil {
		result["created"] = image.Created.Format(time.RFC3339)
	}

	if image.Expiry != nil {
		result["expiry"] = image.Expiry.Format(time.RFC3339)
	}

	return result
}

func imageValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "deprecated", "is_public":
		return strconv.ParseBool(value)

	case "size":
		return strconv.Atoi(value)
	}

	return value, nil
}
