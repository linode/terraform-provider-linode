package linode

import (
	"fmt"
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
				Optional:    true,
				Computed:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeInt,
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
				Description: "If a Volume is attached to a specific Linode, the ID of that Linode will be displayed here.",
				Computed:    true,
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

	_, err = client.GetVolume(int(id))
	if err != nil {
		return false, fmt.Errorf("Failed to get Linode Volume ID %s because %s", d.Id(), err)
	}
	return true, nil
}

func syncVoumeResourceData(d *schema.ResourceData, volume *linodego.Volume) {
	d.Set("label", volume.Label)
	d.Set("region", volume.Region)
	d.Set("status", volume.Status)
	d.Set("size", volume.Size)
	d.Set("linode_id", volume.LinodeID)
	d.Set("filesystem_path", volume.FilesystemPath)
}

func resourceLinodeVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume ID %s as int because %s", d.Id(), err)
	}

	volume, err := client.GetVolume(int(id))

	if err != nil {
		return fmt.Errorf("Failed to find the specified Linode Volume because %s", err)
	}

	syncVoumeResourceData(d, volume)

	return nil
}

func resourceLinodeVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Volume")
	}
	d.Partial(true)

	createOpts := linodego.VolumeCreateOptions{
		Label:  d.Get("label").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(int),
	}
	volume, err := client.CreateVolume(createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode Volume because %s", err)
	}
	d.SetPartial("label")
	d.SetPartial("region")
	d.SetPartial("size")

	if d.HasChange("linode_id") {
		attachOptions := linodego.VolumeAttachOptions{
			LinodeID: d.Get("linode_id").(int),
			ConfigID: 0,
		}
		if ok, err := client.AttachVolume(volume.ID, &attachOptions); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("Failed to attach Linode Volume %d to Linode %d", volume.ID, d.Get("linode_id").(int))
		}
		// @TODO is this the correct use of d.Timeout()?
		linodego.WaitForVolumeStatus(&client, volume.ID, linodego.VolumeActive, int(d.Timeout("update").Seconds()))
		d.Set("linode_id", volume.Label)
		d.SetPartial("linode_id")
	}

	d.SetId(fmt.Sprintf("%d", volume.ID))
	syncVoumeResourceData(d, volume)

	return resourceLinodeVolumeRead(d, meta)
}

func resourceLinodeVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	d.Partial(true)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume id %s as an int because %s", d.Id(), err)
	}

	volume, err := client.GetVolume(int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current linode because %s", err)
	}

	if d.HasChange("size") {
		size := d.Get("size").(int)
		if ok, err := client.ResizeVolume(volume.ID, size); err != nil {
			return err
		} else if ok {
			d.Set("size", size)
			d.SetPartial("size")
		}
	}

	if d.HasChange("label") {
		if volume, err = client.RenameVolume(volume.ID, d.Get("label").(string)); err != nil {
			return err
		}
		d.Set("label", volume.Label)
		d.SetPartial("label")
	}

	if d.HasChange("linode_id") {
		linodeID := d.Get("linode_id").(int)

		if linodeID > 0 {
			attachOptions := linodego.VolumeAttachOptions{
				LinodeID: linodeID,
				ConfigID: 0,
			}
			if ok, err := client.AttachVolume(volume.ID, &attachOptions); err != nil {
				return err
			} else if !ok {
				return fmt.Errorf("Failed to attach Linode Volume %d to Linode %d", volume.ID, linodeID)
			}
			// @TODO is this the correct use of d.Timeout()?
			linodego.WaitForVolumeStatus(&client, volume.ID, linodego.VolumeActive, int(d.Timeout("update").Seconds()))
		} else {
			if ok, err := client.DetachVolume(volume.ID); err != nil {
				return err
			} else if !ok {
				return fmt.Errorf("Failed to detach Linode Volume %d", volume.ID)
			}
			// @TODO do we need to wait for the detach, or should we wait before an attach? what happens when moving a volume to a different label.
			// we'll have to wait for the LinodeID to settle at nil to know that it is unattached
		}
		d.Set("linode_id", volume.ID)
		d.SetPartial("linode_id")
	}

	return nil // resourceLinodeVolumeRead(d, meta)
}

func resourceLinodeVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Volume id %s as int", d.Id())
	}
	err = client.DeleteVolume(int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode Volume %d because %s", id, err)
	}
	return nil
}
