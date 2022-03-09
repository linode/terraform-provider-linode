package databasemysql

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = meta.(*helper.ProviderMeta).Client
	_, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	// TODO: Read logic

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = meta.(*helper.ProviderMeta).Client

	// TODO: Create logic

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = meta.(*helper.ProviderMeta).Client

	_, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	// TODO: Update logic

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_ = meta.(*helper.ProviderMeta).Client
	_, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode MySQL database ID %s as int: %s", d.Id(), err)
	}

	// TODO: Delete logic

	d.SetId("")

	return nil
}
