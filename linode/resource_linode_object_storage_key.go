package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeObjectStorageKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeObjectStorageKeyCreate,
		Read:   resourceLinodeObjectStorageKeyRead,
		Update: resourceLinodeObjectStorageKeyUpdate,
		Delete: resourceLinodeObjectStorageKeyDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label given to this key. For display purposes only.",
				Required:    true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "This keypair's access key. This is not secret.",
				Computed:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "This keypair's secret key.",
				Sensitive:   true,
				Computed:    true,
			},
			"limited": {
				Type:        schema.TypeBool,
				Description: "Whether or not this key is a limited access key.",
				Computed:    true,
			},
			"bucket_access": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket_name": {
							Type:        schema.TypeString,
							Description: "The unique label of the bucket to which the key will grant limited access.",
							Required:    true,
						},
						"cluster": {
							Type:        schema.TypeString,
							Description: "The Object Storage cluster where a bucket to which the key is granting access is hosted.",
							Required:    true,
						},
						"permissions": {
							Type:        schema.TypeString,
							Description: "This Limited Access Key’s permissions for the selected bucket.",
							Required:    true,
						},
					},
				},
				ForceNew: true,
			},
		},
	}
}

func resourceLinodeObjectStorageKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	createOpts := linodego.ObjectStorageKeyCreateOptions{
		Label: d.Get("label").(string),
	}

	if bucketAccess, bucketAccessOk := d.GetOk("bucket_access"); bucketAccessOk {
		createOpts.BucketAccess = expandLinodeObjectStorageKeyBucketAccess(bucketAccess.([]interface{}))
	}

	objectStorageKey, err := client.CreateObjectStorageKey(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode Object Storage Key: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", objectStorageKey.ID))
	d.Set("label", objectStorageKey.Label)
	d.Set("access_key", objectStorageKey.AccessKey)

	// secret_key only available on creation
	d.Set("secret_key", objectStorageKey.SecretKey)

	d.Set("limited", objectStorageKey.Limited)

	bucketAccess := flattenLinodeObjectStorageKeyBucketAccess(objectStorageKey.BucketAccess)
	if bucketAccess != nil {
		d.Set("bucket_access", bucketAccess)
	}

	return resourceLinodeObjectStorageKeyRead(d, meta)
}

func resourceLinodeObjectStorageKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Object Storage Key ID %s as int: %s", d.Id(), err)
	}

	objectStorageKey, err := client.GetObjectStorageKey(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error finding the specified Linode Object Storage Key: %s", err)
	}

	d.Set("label", objectStorageKey.Label)
	d.Set("access_key", objectStorageKey.AccessKey)
	d.Set("limited", objectStorageKey.Limited)

	bucketAccess := flattenLinodeObjectStorageKeyBucketAccess(objectStorageKey.BucketAccess)
	if bucketAccess != nil {
		d.Set("bucket_access", bucketAccess)
	}
	return nil
}

func resourceLinodeObjectStorageKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Object Storage Key id %s as int: %s", d.Id(), err)
	}

	if d.HasChange("label") {
		objectStorageKey, err := client.GetObjectStorageKey(context.Background(), int(id))

		updateOpts := linodego.ObjectStorageKeyUpdateOptions{
			Label: d.Get("label").(string),
		}

		if err != nil {
			return fmt.Errorf("Error fetching data about the current Linode Object Storage Key: %s", err)
		}

		if objectStorageKey, err = client.UpdateObjectStorageKey(context.Background(), int(id), updateOpts); err != nil {
			return err
		}
		d.Set("label", objectStorageKey.Label)
	}

	return resourceLinodeObjectStorageKeyRead(d, meta)
}

func resourceLinodeObjectStorageKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Object Storage Key id %s as int", d.Id())
	}
	err = client.DeleteObjectStorageKey(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode Object Storage Key %d: %s", id, err)
	}
	return nil
}

func flattenLinodeObjectStorageKeyBucketAccess(bucketAccesses *[]linodego.ObjectStorageKeyBucketAccess) *[]map[string]interface{} {
	if bucketAccesses == nil {
		return nil
	}
	specs := make([]map[string]interface{}, len(*bucketAccesses))

	for i, bucketAccess := range *bucketAccesses {
		specs[i] = map[string]interface{}{
			"bucket_name": bucketAccess.BucketName,
			"cluster":     bucketAccess.Cluster,
			"permissions": bucketAccess.Permissions,
		}
	}
	return &specs
}

func expandLinodeObjectStorageKeyBucketAccess(bucketAccessSpecs []interface{}) *[]linodego.ObjectStorageKeyBucketAccess {
	bucketAccesses := make([]linodego.ObjectStorageKeyBucketAccess, len(bucketAccessSpecs))
	for i, bucketAccessSpec := range bucketAccessSpecs {
		bucketAccessSpec := bucketAccessSpec.(map[string]interface{})
		bucketAccess := linodego.ObjectStorageKeyBucketAccess{
			BucketName:  bucketAccessSpec["bucket_name"].(string),
			Cluster:     bucketAccessSpec["cluster"].(string),
			Permissions: bucketAccessSpec["permissions"].(string),
		}
		bucketAccesses[i] = bucketAccess
	}
	return &bucketAccesses
}
