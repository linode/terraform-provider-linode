package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeVolumeCreate,
		Read:   resourceLinodeVolumeRead,
		Update: resourceLinodeVolumeUpdate,
		Delete: resourceLinodeVolumeDelete,
		Exists: resourceLinodeVolumeExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The label of the Linode Volume.",
				Required:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The status of the volume, indicating the current readiness state.",
				Computed:    true,
			},
			"region": &schema.Schema{
				Type:         schema.TypeString,
				Description:  "The region where this volume will be deployed.",
				Required:     true,
				ForceNew:     true,
				InputDefault: "us-east",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Size of the Volume in GB",
				Optional:    true,
				Computed:    true,
			},
			"linode_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The Linode ID where the Volume should be attached.",
				Optional:    true,
			},
			"filesystem_path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The full filesystem path for the Volume based on the Volume's label. Path is /dev/disk/by-id/scsi-0Linode_Volume_ + Volume label.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode Volume ID %s as int because %s", d.Id(), err)
	}

	_, err = client.GetVolume(context.Background(), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Failed to get Linode Volume ID %s because %s", d.Id(), err)
	}
	return true, nil
}

func syncVolumeResourceData(d *schema.ResourceData, volume *linodego.Volume) {
	d.Set("label", volume.Label)
	d.Set("region", volume.Region)
	d.Set("status", volume.Status)
	d.Set("size", volume.Size)
	// d.Set("linode_id", volume.LinodeID)
	d.Set("filesystem_path", volume.FilesystemPath)
}

func resourceLinodeVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume ID %s as int because %s", d.Id(), err)
	}

	volume, err := client.GetVolume(context.Background(), int(id))

	if err != nil {
		if lerr, ok := err.(linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Failed to find the specified Linode Volume because %s", err)
	}

	syncVolumeResourceData(d, volume)

	return nil
}

func resourceLinodeVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Volume")
	}
	d.Partial(true)

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

	volume, err := client.CreateVolume(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode Volume because %s", err)
	}

	d.SetId(fmt.Sprintf("%d", volume.ID))
	d.SetPartial("label")
	d.SetPartial("region")
	d.SetPartial("size")

	if createOpts.LinodeID > 0 {
		if err := linodego.WaitForVolumeLinodeID(context.Background(), &client, volume.ID, linodeID, int(d.Timeout("update").Seconds())); err != nil {
			return err
		}
		d.SetPartial("linode_id")
	}

	syncVolumeResourceData(d, volume)
	if err = linodego.WaitForVolumeStatus(context.Background(), &client, volume.ID, linodego.VolumeActive, int(d.Timeout("create").Seconds())); err != nil {
		return err
	}

	d.Partial(false)
	return resourceLinodeVolumeRead(d, meta)
}

func resourceLinodeVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	d.Partial(true)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume id %s as an int because %s", d.Id(), err)
	}

	volume, err := client.GetVolume(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current linode because %s", err)
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		if ok, err := client.ResizeVolume(context.Background(), volume.ID, size); err != nil {
			return err
		} else if ok {
			if err := linodego.WaitForVolumeStatus(context.Background(), &client, volume.ID, linodego.VolumeActive, int(d.Timeout("update").Seconds())); err != nil {
				return err
			}

			d.Set("size", size)
			d.SetPartial("size")
		}
	}

	if d.HasChange("label") {
		if volume, err = client.RenameVolume(context.Background(), volume.ID, d.Get("label").(string)); err != nil {
			return err
		}
		d.Set("label", volume.Label)
		d.SetPartial("label")
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
			if ok, err := client.DetachVolume(context.Background(), volume.ID); err != nil {
				return err
			} else if !ok {
				return fmt.Errorf("Failed to detach Linode Volume %d", volume.ID)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", volume.ID)
			if err := linodego.WaitForVolumeLinodeID(context.Background(), &client, volume.ID, nil, int(d.Timeout("update").Seconds())); err != nil {
				return err
			}
		}

		if linodeID != nil {
			attachOptions := linodego.VolumeAttachOptions{
				LinodeID: *linodeID,
				ConfigID: 0,
			}

			log.Printf("[INFO] Attaching Linode Volume %d to Linode Instance %d", volume.ID, *linodeID)

			if ok, err := client.AttachVolume(context.Background(), volume.ID, &attachOptions); err != nil {
				return err
			} else if !ok {
				return fmt.Errorf("Failed to attach Linode Volume %d to Linode Instance %d", volume.ID, *linodeID)
			}

			log.Printf("[INFO] Waiting for Linode Volume %d to attach ...", volume.ID)
			if err := linodego.WaitForVolumeLinodeID(context.Background(), &client, volume.ID, linodeID, int(d.Timeout("update").Seconds())); err != nil {
				return err
			}
		}

		d.Set("linode_id", linodeID)
		d.SetPartial("linode_id")
		d.Partial(false)
	}

	return resourceLinodeVolumeRead(d, meta)
}

func resourceLinodeVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id64, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume id %s as int", d.Id())
	}
	id := int(id64)

	log.Printf("[INFO] Detaching Linode Volume %d for deletion", id)
	if ok, err := client.DetachVolume(context.Background(), id); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("Failed to detach Linode Volume %d", id)
	}

	log.Printf("[INFO] Waiting for Linode Volume %d to detach ...", id)
	if err := linodego.WaitForVolumeLinodeID(context.Background(), &client, id, nil, int(d.Timeout("update").Seconds())); err != nil {
		return err
	}

	err = client.DeleteVolume(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode Volume %d because %s", id, err)
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
		// reattach (linode id changed)
		changed = changed || (*have != *want)
	}
	return changed
}
