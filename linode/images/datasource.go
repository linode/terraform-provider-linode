package images

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"context"
	"strconv"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	latestFlag := d.Get("latest").(bool)

	filterID, err := helper.GetFilterID(d)
	if err != nil {
		return diag.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := helper.ConstructFilterString(d, imageValueToFilterType)
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
		imagesFlattened[i] = flattenImage(&image)
	}

	imagesFiltered, err := helper.FilterResults(d, imagesFlattened)
	if err != nil {
		return diag.Errorf("failed to filter returned images: %s", err)
	}

	if latestFlag {
		latestImage := getLatestImage(imagesFiltered)

		if latestImage != nil {
			imagesFiltered = []map[string]interface{}{latestImage}
		}
	}

	d.SetId(filterID)
	d.Set("images", imagesFiltered)

	return nil
}

func flattenImage(image *linodego.Image) map[string]interface{} {
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

func getLatestImage(images []map[string]interface{}) map[string]interface{} {
	var latestCreated time.Time
	var latestImage map[string]interface{}

	for _, image := range images {
		created, ok := image["created"]
		if !ok {
			continue
		}

		createdTime, err := time.Parse(time.RFC3339, created.(string))
		if err != nil {
			return nil
		}

		if latestImage != nil && !createdTime.After(latestCreated) {
			continue
		}

		latestCreated = createdTime
		latestImage = image
	}

	return latestImage
}
