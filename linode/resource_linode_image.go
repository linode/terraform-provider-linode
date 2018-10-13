package linode

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeImageCreate,
		Read:   resourceLinodeImageRead,
		Update: resourceLinodeImageUpdate,
		Delete: resourceLinodeImageDelete,
		Exists: resourceLinodeImageExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode Image.",
				Required:    true,
			},
			"disk_id": {
				Type:        schema.TypeInt,
				Description: "The Disk ID to base the Image on.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A text description of the Image.",
				Optional:    true,
			},
			"image_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Image.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeImageExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)

	_, err := client.GetImage(context.Background(), d.Id())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode Image ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	image, err := client.GetImage(context.Background(), d.Id())

	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Image ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error finding the specified Linode Image: %s", err)
	}

	d.Set("label", image.Label)
	d.Set("description", image.Description)

	return nil
}

func resourceLinodeImageCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Image")
	}
	d.Partial(true)

	createOpts := linodego.ImageCreateOptions{
		DiskID:      d.Get("disk_id").(int),
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}

	image, err := client.CreateImage(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode Image: %s", err)
	}

	d.SetId(image.ID)
	d.SetPartial("label")
	d.SetPartial("description")
	d.Partial(false)

	return resourceLinodeImageRead(d, meta)
}

func resourceLinodeImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	d.Partial(true)

	image, err := client.GetImage(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error fetching data about the current Image: %s", err)
	}

	updateOpts := linodego.ImageUpdateOptions{}

	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)
	}

	if d.HasChange("description") {
		descString := d.Get("description").(string)
		updateOpts.Description = &descString
	}

	image, err = client.UpdateImage(context.Background(), d.Id(), updateOpts)
	if err != nil {
		return err
	}

	d.Set("label", image.Label)
	d.Set("description", image.Description)

	return resourceLinodeImageRead(d, meta)
}

func resourceLinodeImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	log.Printf("[INFO] Detaching Linode Image %s for deletion", d.Id())

	err := client.DeleteImage(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Linode Image %s: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
