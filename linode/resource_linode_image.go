package linode

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

const (
	LinodeImageCreateTimeout = 20 * time.Minute
)

func resourceLinodeImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeImageCreate,
		ReadContext:   resourceLinodeImageRead,
		UpdateContext: resourceLinodeImageUpdate,
		DeleteContext: resourceLinodeImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(LinodeImageCreateTimeout),
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "A short description of the Image. Labels cannot contain special characters.",
				Required:    true,
			},
			"disk_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the Linode Disk that this Image will be created from.",
				RequiredWith:  []string{"linode_id"},
				ConflictsWith: []string{"file_path"},
				Optional:      true,
				ForceNew:      true,
			},
			"linode_id": {
				Type:          schema.TypeInt,
				Description:   "The ID of the Linode that this Image will be created from.",
				RequiredWith:  []string{"disk_id"},
				ConflictsWith: []string{"file_path"},
				Optional:      true,
				ForceNew:      true,
			},
			"file_path": {
				Type:          schema.TypeString,
				Description:   "The name of the file to upload to this image.",
				ConflictsWith: []string{"linode_id", "disk_id"},
				RequiredWith:  []string{"region"},
				Optional:      true,
				ForceNew:      true,
			},
			"region": {
				Type:         schema.TypeString,
				Description:  "The region to upload to.",
				RequiredWith: []string{"file_path"},
				Optional:     true,
			},
			"file_hash": {
				Type:        schema.TypeString,
				Description: "The MD5 hash of the image file.",
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A detailed description of this Image.",
				Optional:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "When this Image was created.",
				Computed:    true,
			},
			"created_by": {
				Type:        schema.TypeString,
				Description: "The name of the User who created this Image.",
				Computed:    true,
			},
			"deprecated": {
				Type:        schema.TypeBool,
				Description: "Whether or not this Image is deprecated. Will only be True for deprecated public Images.",
				Computed:    true,
			},
			"is_public": {
				Type:        schema.TypeBool,
				Description: "True if the Image is public.",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "The minimum size this Image needs to deploy. Size is in MB.",
				Computed:    true,
			},
			"type": {
				Type: schema.TypeString,
				Description: "How the Image was created. 'Manual' Images can be created at any time. 'Automatic' " +
					"images are created automatically from a deleted Linode.",
				Computed: true,
			},
			"expiry": {
				Type:        schema.TypeString,
				Description: "Only Images created automatically (from a deleted Linode; type=automatic) will expire.",
				Computed:    true,
			},
			"vendor": {
				Type:        schema.TypeString,
				Description: "The upstream distribution vendor. Nil for private Images.",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "The current status of this Image.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	image, err := client.GetImage(ctx, d.Id())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode Image ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error getting Linode image %s: %s", d.Id(), err)
	}

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

func resourceLinodeImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if _, ok := d.GetOk("linode_id"); ok {
		return resourceLinodeImageCreateFromLinode(ctx, d, meta)
	}

	if _, ok := d.GetOk("file_path"); ok {
		return resourceLinodeImageCreateFromUpload(ctx, d, meta)
	}

	return diag.Errorf("failed to create image: source or linode_id must be specified")
}

func resourceLinodeImageCreateFromLinode(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)
	diskID := d.Get("disk_id").(int)

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, int(d.Timeout(schema.TimeoutCreate).Seconds()),
	); err != nil {
		return diag.Errorf(
			"Error waiting for Linode Instance %d Disk %d to become ready for taking an Image", linodeID, diskID)
	}

	createOpts := linodego.ImageCreateOptions{
		DiskID:      diskID,
		Label:       d.Get("label").(string),
		Description: d.Get("description").(string),
	}

	image, err := client.CreateImage(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Image: %s", err)
	}

	d.SetId(image.ID)

	if _, err := client.WaitForInstanceDiskStatus(
		ctx, linodeID, diskID, linodego.DiskReady, int(d.Timeout(schema.TimeoutCreate).Seconds()),
	); err != nil {
		return diag.Errorf(
			"failed to wait for linode instance %d disk %d to become ready while taking an image", linodeID, diskID)
	}

	return resourceLinodeImageRead(ctx, d, meta)
}

func resourceLinodeImageCreateFromUpload(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	region := d.Get("region").(string)
	label := d.Get("label").(string)
	description := d.Get("description").(string)

	imageReader, err := imageFromResourceData(d)
	if err != nil {
		diag.Errorf("failed to get image source: %v", err)
	}
	defer imageReader.Close()

	createOpts := linodego.ImageCreateUploadOptions{
		Region:      region,
		Label:       label,
		Description: description,
	}

	image, uploadURL, err := client.CreateImageUpload(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create image upload %s: %v", label, err)
	}

	if err := uploadImageAndStoreHash(ctx, d, meta, uploadURL, imageReader); err != nil {
		return diag.Errorf("failed to upload image: %v", err)
	}

	image, err = client.WaitForImageStatus(ctx, image.ID, linodego.ImageStatusAvailable,
		int(d.Timeout(schema.TimeoutCreate).Seconds()))
	if err != nil {
		return diag.Errorf("failed to wait for image to be available: %v", err)
	}

	d.SetId(image.ID)

	return resourceLinodeImageRead(ctx, d, meta)
}

func resourceLinodeImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

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
		image, err = client.UpdateImage(ctx, d.Id(), updateOpts)
		if err != nil {
			return diag.Errorf("failed to update image: %v", err)
		}
	}

	d.Set("label", image.Label)
	d.Set("description", image.Description)

	return resourceLinodeImageRead(ctx, d, meta)
}

func resourceLinodeImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

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
	uploadURL string, image io.Reader) error {
	client := meta.(*ProviderMeta).Client

	var buf bytes.Buffer
	tee := io.TeeReader(image, &buf)

	if err := client.UploadImageToURL(ctx, uploadURL, tee); err != nil {
		return err
	}

	hash := md5.New()

	if _, err := io.Copy(hash, &buf); err != nil {
		return err
	}

	d.Set("file_hash", hex.EncodeToString(hash.Sum(nil)))

	return nil
}

func flattenLinodeImage(image *linodego.Image) map[string]interface{} {
	result := make(map[string]interface{})

	result["id"] = image.ID
	result["label"] = image.Label
	result["description"] = image.Description
	result["created_by"] = image.CreatedBy
	result["deprecated"] = image.Deprecated
	result["is_public"] = image.IsPublic
	result["size"] = image.Size
	result["type"] = image.Type
	result["vendor"] = image.Vendor
	result["status"] = image.Status

	if image.Created != nil {
		result["created"] = image.Created.Format(time.RFC3339)
	}

	if image.Expiry != nil {
		result["expiry"] = image.Expiry.Format(time.RFC3339)
	}

	return result
}
