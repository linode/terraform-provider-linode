package rdns

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const updateRDNSTimeout = time.Minute * 10

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		UpdateContext: updateResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(updateRDNSTimeout),
			Update: schema.DefaultTimeout(updateRDNSTimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	ip, err := client.GetIPAddress(ctx, ipStr)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode RDNS %q from state because it no longer exists", ipStr)
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode RDNS: %s", err)
	}

	d.Set("address", d.Id())
	d.Set("rdns", ip.RDNS)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	address := d.Get("address").(string)
	var rdns *string
	if rdnsRaw, ok := d.GetOk("rdns"); ok && len(rdnsRaw.(string)) > 0 {
		rdnsStr := rdnsRaw.(string)
		rdns = &rdnsStr
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: rdns,
	}

	ip, err := updateIPAddress(ctx, d, meta, address, updateOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode RDNS: %s", err)
	}
	d.SetId(address)
	d.Set("rdns", ip.RDNS)

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	var rdns *string

	if rdnsRaw, ok := d.GetOk("rdns"); ok && len(rdnsRaw.(string)) > 0 {
		rdnsStr := rdnsRaw.(string)
		rdns = &rdnsStr
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: rdns,
	}

	if _, err := updateIPAddress(ctx, d, meta, d.Id(), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode RDNS: %s", err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: nil,
	}

	if _, err := client.UpdateIPAddress(ctx, d.Id(), updateOpts); err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error deleting Linode RDNS: %s", err)
	}

	d.SetId("")

	return nil
}

// updateIPAddress wraps the client.UpdateIPAddress(...) retry logic depending on the 'wait_for_available' field.
func updateIPAddress(ctx context.Context, d *schema.ResourceData, meta interface{}, address string,
	updateOpts linodego.IPAddressUpdateOptions) (*linodego.InstanceIP, error) {
	client := meta.(*helper.ProviderMeta).Client
	retry := d.Get("wait_for_available").(bool)

	if retry {
		return updateIPAddressWithRetries(ctx, &client, address, updateOpts, time.Second)
	}

	return client.UpdateIPAddress(ctx, address, updateOpts)
}

func updateIPAddressWithRetries(ctx context.Context, client *linodego.Client, address string,
	updateOpts linodego.IPAddressUpdateOptions, retryDuration time.Duration) (*linodego.InstanceIP, error) {
	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := client.UpdateIPAddress(ctx, address, updateOpts)
			if err != nil {
				if lerr, ok := err.(*linodego.Error); ok && lerr.Code != 400 &&
					!strings.Contains(lerr.Error(), "unable to perform a lookup") {
					return nil, fmt.Errorf("failed to update ip address: %s", err)
				}

				continue
			}

			return result, nil

		case <-ctx.Done():
			return nil, fmt.Errorf("failed to update ip address: %s", ctx.Err())
		}
	}
}
