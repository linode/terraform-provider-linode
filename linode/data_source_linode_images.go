package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"

	"context"
	"fmt"
	"strconv"
)

func dataSourceLinodeImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeImagesRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema([]string{"deprecated", "is_public", "label", "size", "vendor"}),
			"images": {
				Type:        schema.TypeList,
				Description: "The returned list of Images.",
				Computed:    true,
				Elem:        dataSourceLinodeImage(),
			},
		},
	}
}

func dataSourceLinodeImagesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	filter, err := constructFilterString(d, imageValueToFilterType)
	if err != nil {
		return fmt.Errorf("failed to construct filter: %s", err)
	}

	images, err := client.ListImages(context.Background(), &linodego.ListOptions{
		Filter: filter,
	})

	if err != nil {
		return fmt.Errorf("failed to list linode images: %s", err)
	}

	imagesFlattened := make([]interface{}, len(images))
	for i, image := range images {
		imagesFlattened[i] = flattenLinodeImage(&image)
	}

	d.SetId(filter)
	d.Set("images", imagesFlattened)

	return nil
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
