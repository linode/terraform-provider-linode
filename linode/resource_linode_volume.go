package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	LinodeVolumeCreateTimeout = 10 * time.Minute
	LinodeVolumeUpdateTimeout = 20 * time.Minute
	LinodeVolumeDeleteTimeout = 10 * time.Minute
)

func resourceLinodeVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeVolumeCreateContext,
		ReadContext:   resourceLinodeVolumeReadContext,
		UpdateContext: resourceLinodeVolumeUpdateContext,
		DeleteContext: resourceLinodeVolumeDeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LinodeVolumeCreateTimeout),
			Update: schema.DefaultTimeout(LinodeVolumeUpdateTimeout),
			Delete: schema.DefaultTimeout(LinodeVolumeDeleteTimeout),
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode Volume.",
				Required:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The status of the volume, indicating the current readiness state.",
				Computed:    true,
			},
			"region": {
				Type:         schema.TypeString,
				Description:  "The region where this volume will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Size of the Volume in GB",
				Optional:    true,
				Computed:    true,
			},
			"linode_id": {
				Type:        schema.TypeInt,
				Description: "The Linode ID where the Volume should be attached.",
				Optional:    true,
				Computed:    true,
			},
			"filesystem_path": {
				Type:        schema.TypeString,
				Description: "The full filesystem path for the Volume based on the Volume's label. Path is /dev/disk/by-id/scsi-0Linode_Volume_ + Volume label.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			},
		},
	}
}

func resourceLinodeVolumeReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume ID %s as int: %s", d.Id(), err)
	}

	volume, err := client.GetVolume(context.Background(), int(id))

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

func resourceLinodeVolumeCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(linodego.Client)
	if !ok {
		return diag.Errorf("Invalid Client when creating Linode Volume")
	}

	var linodeID *int

	createOpts := linodego.VolumeCreateOptions{
		Label:  d.Get("label").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(int),
	}

	if lID, ok := d.GetOk("linode_id"); ok {
		lidInt := lID.(int)
		linodeID = &lidInt
		createOpts.LinodeID = *linodeID
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Volume: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", volume.ID))

	if createOpts.LinodeID > 0 {
		if _, err := client.WaitForVolumeLinodeID(context.Background(), volume.ID, linodeID, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, int(d.Timeout(schema.TimeoutCreate).Seconds())); err != nil {
		return diag.FromErr(err)
	}

	return resourceLinodeVolumeReadContext(ctx, d, meta)
}

func resourceLinodeVolumeUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume id %s as int: %s", d.Id(), err)
	}

	volume, errVolume := client.GetVolume(context.Background(), int(id))
	if errVolume != nil {
		return diag.Errorf("Error fetching data about the volume %d: %s", int(id), errVolume)
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		if err = client.ResizeVolume(context.Background(), volume.ID, size); err != nil {
			return diag.FromErr(err)
		}

		if _, err = client.WaitForVolumeStatus(context.Background(), volume.ID, linodego.VolumeActive, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
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
		if volume, err = client.UpdateVolume(context.Background(), volume.ID, updateOpts); err != nil {
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

	if detectVolumeIDChange(linodeID, volume.LinodeID) {
		if linodeID == nil || volume.LinodeID != nil {
			log.Printf("[INFO] Detaching Linode Volume %d", volume.ID)
			if err = client.DetachVolume(context.Background(), volume.ID); err != nil {
				return diag.FromErr(err)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", volume.ID)
			if _, err = client.WaitForVolumeLinodeID(context.Background(), volume.ID, nil, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
				return diag.FromErr(err)
			}
		}

		if linodeID != nil {
			attachOptions := linodego.VolumeAttachOptions{
				LinodeID: *linodeID,
				ConfigID: 0,
			}

			log.Printf("[INFO] Attaching Linode Volume %d to Linode Instance %d", volume.ID, *linodeID)

			if _, err = client.AttachVolume(context.Background(), volume.ID, &attachOptions); err != nil {
				return diag.Errorf("Error attaching Linode Volume %d to Linode Instance %d: %s", volume.ID, *linodeID, err)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to attach ...", volume.ID)
			if _, err = client.WaitForVolumeLinodeID(context.Background(), volume.ID, linodeID, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
				return diag.FromErr(err)
			}
		}

		d.Set("linode_id", linodeID)
	}

	return resourceLinodeVolumeReadContext(ctx, d, meta)
}

func resourceLinodeVolumeDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id64, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Volume id %s as int", d.Id())
	}
	id := int(id64)

	log.Printf("[INFO] Detaching Linode Volume %d for deletion", id)
	if err := client.DetachVolume(context.Background(), id); err != nil {
		return diag.Errorf("Error detaching Linode Volume %d: %s", id, err)
	}

	log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", id)
	if _, err := client.WaitForVolumeLinodeID(context.Background(), id, nil, int(d.Timeout(schema.TimeoutUpdate).Seconds())); err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteVolume(context.Background(), int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode Volume %d: %s", id, err)
	}
	d.SetId("")
	return nil
}

func detectVolumeIDChange(have *int, want *int) (changed bool) {
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
