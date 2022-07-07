package images

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listImages, flattenImage)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("images", results)

	return nil
}

func listImages(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client,
	options *linodego.ListOptions,
) ([]interface{}, error) {
	images, err := client.ListImages(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(images))

	for i, v := range images {
		result[i] = v
	}

	return result, nil
}

func flattenImage(data interface{}) map[string]interface{} {
	image := data.(linodego.Image)

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
