package image

import (
	"bytes"
	"context"
	"crypto/md5" // #nosec G501 -- endpoint expecting md5
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	LinodeImageCreateTimeout = 30 * time.Minute
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
			Create: schema.DefaultTimeout(LinodeImageCreateTimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_image")

	client := meta.(*helper.ProviderMeta).Client

	tflog.Trace(ctx, "client.GetImage(...)")
	image, err := client.GetImage(ctx, d.Id())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode Image ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error getting Linode image %s: %s", d.Id(), err)
	}

	d.Set("capabilities", image.Capabilities)
	d.Set("label", image.Label)
	d.Set("description", image.Description)
	d.Set("type", image.Type)
	d.Set("size", image.Size)
	d.Set("vendor", image.Vendor)
	d.Set("created_by", image.CreatedBy)
	d.Set("deprecated", image.Deprecated)
	d.Set("is_public", image.IsPublic)
	d.Set("status", image.Status)

	if image.Created != nil {
		d.Set("created", image.Created.Format(time.RFC3339))
	}
	if image.Expiry != nil {
		d.Set("expiry", image.Expiry.Format(time.RFC3339))
	}

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("linode_id"); ok {
		return createResourceFromLinode(ctx, d, meta)
	}

	if _, ok := d.GetOk("file_path"); ok {
		return createResourceFromUpload(ctx, d, meta)
	}

	return diag.Errorf("failed to create image: source or linode_id must be specified")
}

func createResourceFromLinode(
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
	tflog.Debug(ctx, "Create linode_image")

	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)
	diskID := d.Get("disk_id").(int)

	timeoutSeconds, err := helper.SafeFloat64ToInt(
		d.Timeout(schema.TimeoutCreate).Seconds(),
	)
	if err != nil {
		return diag.Errorf(
			"failed to convert the float64 creation timeout to be an integer: %s",
			err,
		)
	}

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, timeoutSeconds,
	); err != nil {
		return diag.Errorf(
			"Error waiting for Linode Instance %d Disk %d to become ready for taking an Image", linodeID, diskID)
	}

	createOpts := linodego.ImageCreateOptions{
		DiskID:      diskID,
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
		CloudInit:   d.Get("cloud_init").(bool),
	}

	tflog.Trace(ctx, "client.CreateImage(...)", map[string]any{
		"options": createOpts,
	})

	image, err := client.CreateImage(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Image: %s", err)
	}

	d.SetId(image.ID)

	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Waiting for a single image to be ready")

	tflog.Trace(ctx, "client.WaitForInstanceDiskStatus(...)")

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, timeoutSeconds,
	); err != nil {
		return diag.Errorf(
			"failed to wait for linode instance %d disk %d to become ready while taking an image", linodeID, diskID)
	}

	return readResource(ctx, d, meta)
}

func createResourceFromUpload(
	ctx context.Context, d *schema.ResourceData, meta interface{},
) diag.Diagnostics {
	tflog.Debug(ctx, "Create linode_image_from_upload")

	client := meta.(*helper.ProviderMeta).Client

	region := d.Get("region").(string)
	label := d.Get("label").(string)
	description := d.Get("description").(string)
	cloudInit := d.Get("cloud_init").(bool)

	imageReader, err := imageFromResourceData(d)
	if err != nil {
		return diag.Errorf("failed to get image source: %v", err)
	}

	defer func() {
		if err := imageReader.Close(); err != nil {
			log.Printf("[WARN] Failed to close image reader: %s\n", err)
		}
	}()

	createOpts := linodego.ImageCreateUploadOptions{
		Region:      region,
		Label:       label,
		Description: description,
		CloudInit:   cloudInit,
	}

	tflog.Trace(ctx, "client.CreateImageUpload(...)", map[string]any{
		"options": createOpts,
	})

	image, uploadURL, err := client.CreateImageUpload(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create image upload %s: %v", label, err)
	}

	if err := uploadImageAndStoreHash(ctx, d, meta, uploadURL, imageReader); err != nil {
		return diag.Errorf("failed to upload image: %v", err)
	}
	timeoutSeconds, err := helper.SafeFloat64ToInt(
		d.Timeout(schema.TimeoutCreate).Seconds(),
	)
	if err != nil {
		return diag.Errorf(
			"failed to convert the float64 creation timeout to be an integer: %s",
			err,
		)
	}

	ctx = tflog.SetField(ctx, "image_id", image.ID)
	tflog.Debug(ctx, "Waiting for a single image to be ready")

	tflog.Trace(ctx, "client.WaitForImageStatus(...)")

	image, err = client.WaitForImageStatus(
		ctx,
		image.ID,
		linodego.ImageStatusAvailable,
		timeoutSeconds,
	)
	if err != nil {
		return diag.Errorf("failed to wait for image to be available: %v", err)
	}

	d.SetId(image.ID)

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_image")

	client := meta.(*helper.ProviderMeta).Client

	image, err := client.GetImage(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error fetching data about the current Image: %s", err)
	}

	updateOpts := linodego.ImageUpdateOptions{}

	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)
	}

	if d.HasChange("description") {
		descString := d.Get("description").(string)
		updateOpts.Description = &descString
	}

	if d.HasChanges("label", "description") {
		tflog.Debug(ctx, "client.UpdateImage(...)", map[string]any{
			"options": updateOpts,
		})

		image, err = client.UpdateImage(ctx, d.Id(), updateOpts)
		if err != nil {
			return diag.Errorf("failed to update image: %v", err)
		}
	}

	d.Set("label", image.Label)
	d.Set("description", image.Description)

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_image")

	client := meta.(*helper.ProviderMeta).Client

	tflog.Trace(ctx, "client.DeleteImage(...)")

	err := client.DeleteImage(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting Linode Image %s: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func imageFromResourceData(d *schema.ResourceData) (image io.ReadCloser, err error) {
	if imageFile, ok := d.GetOk("file_path"); ok {
		file, err := os.Open(imageFile.(string))
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", imageFile, err)
		}

		return file, nil
	}

	return nil, fmt.Errorf("no image source specified")
}

func uploadImageAndStoreHash(
	ctx context.Context, d *schema.ResourceData, meta interface{},
	uploadURL string, image io.Reader,
) error {
	client := meta.(*helper.ProviderMeta).Client

	var buf bytes.Buffer
	tee := io.TeeReader(image, &buf)

	if err := client.UploadImageToURL(ctx, uploadURL, tee); err != nil {
		return err
	}

	hash := md5.New() // #nosec G401 -- endpoint expecting md5

	if _, err := io.Copy(hash, &buf); err != nil {
		return err
	}

	d.Set("file_hash", hex.EncodeToString(hash.Sum(nil)))

	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"image_id": d.Id(),
	})
}
