package instancedisk

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"log"
	"strconv"
	"time"
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

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	linodeID := d.Get("linode_id").(int)

	disk, err := client.GetInstanceDisk(ctx, linodeID, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Instance Disk ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Instance Disk: %s", err)
	}

	d.Set("created", disk.Created.Format(time.RFC3339))
	d.Set("filesystem", disk.Filesystem)
	d.Set("label", disk.Label)
	d.Set("size", disk.Size)
	d.Set("status", disk.Status)
	d.Set("updated", disk.Updated.Format(time.RFC3339))

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	createOpts := linodego.InstanceDiskCreateOptions{
		AuthorizedKeys:  helper.ExpandStringSet(d.Get("authorized_keys").(*schema.Set)),
		AuthorizedUsers: helper.ExpandStringSet(d.Get("authorized_users").(*schema.Set)),
		Filesystem:      d.Get("filesystem").(string),
		Image:           d.Get("image").(string),
		Label:           d.Get("label").(string),
		RootPass:        d.Get("root_pass").(string),
		Size:            d.Get("size").(int),
		StackscriptID:   d.Get("stackscript_id").(int),
	}

	if stackscriptData, ok := d.GetOk("stackscript_data"); ok {
		createOpts.StackscriptData = expandStackScriptData(stackscriptData)
	}

	cfg, err := client.CreateInstanceDisk(ctx, linodeID, createOpts)
	if err != nil {
		return diag.Errorf("failed to create linode instance disk: %s", err)
	}

	d.SetId(strconv.Itoa(cfg.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	linodeID := d.Get("linode_id").(int)

	if d.HasChange("size") {
		if err := client.ResizeInstanceDisk(ctx, linodeID, int(id), d.Get("size").(int)); err != nil {
			return diag.Errorf("failed to resize disk: %s", err)
		}
	}

	putRequest := linodego.InstanceDiskUpdateOptions{}
	shouldUpdate := false

	if d.HasChange("label") {
		putRequest.Label = d.Get("label").(string)
		shouldUpdate = true
	}

	if shouldUpdate {
		if _, err := client.UpdateInstanceDisk(ctx, linodeID, int(id), putRequest); err != nil {
			return diag.Errorf("failed to update instance disk: %s", err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id64, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance id %s as int", d.Id())
	}

	linodeID := d.Get("linode_id").(int)

	//disk, err := client.GetInstanceDisk(ctx, linodeID, int(id64))
	//if err != nil {
	//	return diag.Errorf("failed to get disk: %s", err)
	//}

	err = client.DeleteInstanceDisk(ctx, linodeID, int(id64))
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance Disk %d: %s", id64, err)
	}

	d.SetId("")
	return nil
}

func expandStackScriptData(data any) map[string]string {
	dataMap := data.(map[string]any)
	result := make(map[string]string, len(dataMap))

	for k, v := range dataMap {
		result[k] = v.(string)
	}

	return result
}
