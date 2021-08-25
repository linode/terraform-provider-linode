package volume

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
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
	requestedVolumeID := d.Get("id").(int)

	if requestedVolumeID == 0 {
		return diag.Errorf("Volume ID is required")
	}

	var volume *linodego.Volume

	volume, err := client.GetVolume(ctx, requestedVolumeID)
	if err != nil {
		return diag.Errorf("Error requesting Volume: %s", err)
	}

	if volume != nil {
		d.SetId(strconv.Itoa(volume.ID))
		d.Set("region", volume.Region)
		d.Set("size", volume.Size)
		d.Set("filesystem_path", volume.FilesystemPath)
		d.Set("label", volume.Label)
		d.Set("linode_id", volume.LinodeID)
		d.Set("status", volume.Status)
		d.Set("created", volume.Created.Format(time.RFC3339))
		d.Set("updated", volume.Updated.Format(time.RFC3339))
		d.Set("tags", volume.Tags)
		return nil
	}

	return diag.Errorf("Linode Volume %s was not found", fmt.Sprint(requestedVolumeID))
}
