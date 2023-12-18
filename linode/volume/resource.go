package volume

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	LinodeVolumeCreateTimeout = 15 * time.Minute
	LinodeVolumeUpdateTimeout = 20 * time.Minute
	LinodeVolumeDeleteTimeout = 10 * time.Minute
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LinodeVolumeCreateTimeout),
			Update: schema.DefaultTimeout(LinodeVolumeUpdateTimeout),
			Delete: schema.DefaultTimeout(LinodeVolumeDeleteTimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume ID %s as int: %s", d.Id(), err)
	}

	volume, err := client.GetVolume(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Volume ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Volume: %s", err)
	}

	d.Set("label", volume.Label)
	d.Set("region", volume.Region)
	d.Set("status", volume.Status)
	d.Set("size", volume.Size)
	d.Set("linode_id", volume.LinodeID)
	d.Set("filesystem_path", volume.FilesystemPath)
	d.Set("tags", volume.Tags)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	var volume *linodego.Volume
	// Clone the source volume if `source_volume_id` is specified
	if _, ok := d.GetOk("source_volume_id"); ok {
		newVolume, err := createVolumeFromSource(ctx, d, client)
		if err != nil {
			return diag.Errorf("failed to clone volume: %s", err)
		}

		volume = newVolume
	} else {
		newVolume, err := createVolume(ctx, d, client)
		if err != nil {
			return diag.Errorf("failed to create volume: %s", err)
		}

		volume = newVolume
	}

	// Wait for the volume to be created
	if lID, ok := d.GetOk("linode_id"); ok {
		id := lID.(int)

		if _, err := client.WaitForVolumeLinodeID(
			ctx, volume.ID, &id, int(d.Timeout(schema.TimeoutUpdate).Seconds()),
		); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, err := client.WaitForVolumeStatus(
		ctx, volume.ID, linodego.VolumeActive, int(d.Timeout(schema.TimeoutCreate).Seconds()),
	); err != nil {
		return diag.FromErr(err)
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume id %s as int: %s", d.Id(), err)
	}

	volume, errVolume := client.GetVolume(ctx, id)
	if errVolume != nil {
		return diag.Errorf("Error fetching data about the volume %d: %s", id, errVolume)
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		if err = client.ResizeVolume(ctx, volume.ID, size); err != nil {
			return diag.FromErr(err)
		}

		if _, err = client.WaitForVolumeStatus(
			ctx, volume.ID, linodego.VolumeActive, int(d.Timeout(schema.TimeoutUpdate).Seconds()),
		); err != nil {
			return diag.FromErr(err)
		}

		d.Set("size", size)
	}

	updateOpts := linodego.VolumeUpdateOptions{}
	doUpdate := false
	if d.HasChange("tags") {
		tags := []string{}
		for _, tag := range d.Get("tags").(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}

		updateOpts.Tags = &tags
		doUpdate = true
	}

	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)
		doUpdate = true
	}

	if doUpdate {
		if volume, err = client.UpdateVolume(ctx, volume.ID, updateOpts); err != nil {
			return diag.FromErr(err)
		}
		d.Set("tags", volume.Tags)
		d.Set("label", volume.Label)
	}

	var linodeID *int

	if lID, ok := d.GetOk("linode_id"); ok {
		lidInt := lID.(int)
		linodeID = &lidInt
	}

	// We can't use d.HasChange("linode_id") - see https://github.com/hashicorp/terraform/pull/1445
	// compare nils to ints cautiously

	if DetectVolumeIDChange(linodeID, volume.LinodeID) {
		if linodeID == nil || volume.LinodeID != nil {
			log.Printf("[INFO] Detaching Linode Volume %d", volume.ID)
			if err = client.DetachVolume(ctx, volume.ID); err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", volume.ID)
			if _, err = client.WaitForVolumeLinodeID(
				ctx, volume.ID, nil, int(d.Timeout(schema.TimeoutUpdate).Seconds()),
			); err != nil {
				return diag.FromErr(err)
			}
		}

		if linodeID != nil {
			attachOptions := linodego.VolumeAttachOptions{
				LinodeID: *linodeID,
				ConfigID: 0,
			}

			log.Printf("[INFO] Attaching Linode Volume %d to Linode Instance %d", volume.ID, *linodeID)

			if _, err = client.AttachVolume(ctx, volume.ID, &attachOptions); err != nil {
				return diag.Errorf("Error attaching Linode Volume %d to Linode Instance %d: %s", volume.ID, *linodeID, err)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to attach ...", volume.ID)
			if _, err = client.WaitForVolumeLinodeID(
				ctx, volume.ID, linodeID, int(d.Timeout(schema.TimeoutUpdate).Seconds()),
			); err != nil {
				return diag.FromErr(err)
			}
		}

		d.Set("linode_id", linodeID)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume id %s as int", d.Id())
	}

	// We should retry on intermittent deletion errors
	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		vol, err := client.GetVolume(ctx, id)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		// We only want to detach if the volume is already attached to an instance,
		// otherwise a context timeout will be thrown for non-attached volumes.
		if vol.LinodeID != nil && *vol.LinodeID != 0 {
			p, err := client.NewEventPoller(ctx, id, "volume", linodego.ActionVolumeDetach)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("failed to initialize event poller: %s", err))
			}

			log.Printf("[INFO] Detaching Linode Volume %d for deletion\n", id)

			if err := client.DetachVolume(ctx, id); err != nil {
				return retry.RetryableError(fmt.Errorf("error detaching Linode Volume %d: %s", id, err))
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to detach...\n", id)

			if _, err := p.WaitForFinished(ctx, 5); err != nil {
				return retry.RetryableError(err)
			}
		}

		// Let's make sure the volume is fully offline
		if _, err := client.WaitForVolumeLinodeID(
			ctx, id, nil, int(d.Timeout(schema.TimeoutDelete).Seconds()),
		); err != nil {
			return retry.NonRetryableError(err)
		}

		err = client.DeleteVolume(ctx, id)
		if err != nil {
			return retry.RetryableError(fmt.Errorf("error deleting Linode Volume %d: %s", id, err))
		}
		d.SetId("")
		return nil
	}))
}

func DetectVolumeIDChange(have *int, want *int) (changed bool) {
	if have == nil && want == nil {
		changed = false
	} else {
		// attach or detach
		changed = (have == nil && want != nil) || (have != nil && want == nil)
		// reattach (Linode Instance ID changed)
		changed = changed || (*have != *want)
	}
	return changed
}

func createVolume(ctx context.Context, d *schema.ResourceData, client linodego.Client) (*linodego.Volume, error) {
	createOpts := linodego.VolumeCreateOptions{
		Label:  d.Get("label").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(int),
	}

	if lID, ok := d.GetOk("linode_id"); ok {
		lidInt := lID.(int)
		createOpts.LinodeID = lidInt
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	newVolume, err := client.CreateVolume(ctx, createOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create linode volume: %s", err)
	}

	d.SetId(strconv.Itoa(newVolume.ID))

	return newVolume, nil
}

func createVolumeFromSource(
	ctx context.Context, d *schema.ResourceData, client linodego.Client,
) (*linodego.Volume, error) {
	var clonedVolume *linodego.Volume

	newRegion, regionOk := d.GetOk("region")
	newLabel := d.Get("label").(string)
	newSize, sizeOk := d.GetOk("size")

	sourceVolumeID := d.Get("source_volume_id").(int)

	sourceVolume, err := client.GetVolume(ctx, sourceVolumeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source volume %d: %s", sourceVolumeID, err)
	}

	if regionOk && sourceVolume.Region != newRegion.(string) {
		return nil, fmt.Errorf("`region` of source volume differs from specified region: %s != %s",
			sourceVolume.Region, newRegion)
	}

	if sizeOk && sourceVolume.Size > newSize.(int) {
		return nil, fmt.Errorf("`size` must be greater than or equal to the size of the source volume: %d < %d",
			newSize, sourceVolume.Size)
	}

	clonedVolume, err = client.CloneVolume(ctx, sourceVolumeID, newLabel)
	if err != nil {
		return nil, fmt.Errorf("failed to clone volume %d: %s", sourceVolumeID, err)
	}

	d.SetId(strconv.Itoa(clonedVolume.ID))

	if _, err := client.WaitForVolumeStatus(ctx, clonedVolume.ID, linodego.VolumeActive,
		int(LinodeVolumeUpdateTimeout.Seconds())); err != nil {
		return nil, fmt.Errorf("failed to wait for volume to clone: %s", err)
	}

	// Since a cloned volume will have the attributes of the source volume, we need to update
	// to match the schema.
	var updateOpts linodego.VolumeUpdateOptions

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		tags := helper.ExpandStringSet(tagsRaw.(*schema.Set))
		updateOpts.Tags = &tags
	}

	if clonedVolume, err = client.UpdateVolume(ctx, clonedVolume.ID, updateOpts); err != nil {
		return nil, fmt.Errorf("failed to update cloned volume: %s", err)
	}

	// Resize the volume if necessary
	if sizeOk && clonedVolume.Size != newSize.(int) {
		if err := client.ResizeVolume(ctx, clonedVolume.ID, newSize.(int)); err != nil {
			return nil, fmt.Errorf("failed to resize cloned volume %d: %s", clonedVolume.ID, err)
		}

		if _, err := client.WaitForVolumeStatus(
			ctx, clonedVolume.ID, linodego.VolumeActive, int(d.Timeout(schema.TimeoutUpdate).Seconds()),
		); err != nil {
			return nil, fmt.Errorf("failed to wait for volume active: %s", err)
		}
	}

	// Attach the volume if necessary
	if lID, ok := d.GetOk("linode_id"); ok {
		if clonedVolume, err = client.AttachVolume(ctx, clonedVolume.ID, &linodego.VolumeAttachOptions{
			LinodeID: lID.(int),
		}); err != nil {
			return nil, fmt.Errorf("failed to attach cloned volume: %s", err)
		}
	}

	return clonedVolume, nil
}
