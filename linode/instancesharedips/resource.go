package instancesharedips

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_instance_shared_ips")

	helper.AttemptWarnEarlyAccessSDKv2(meta.(*helper.ProviderMeta))

	client := meta.(*helper.ProviderMeta).Client
	linodeID := d.Get("linode_id").(int)

	sharedIPs, err := GetSharedIPsForLinode(ctx, client, linodeID)
	if err != nil {
		return diag.Errorf("failed to get shared ips for linode %d: %s", linodeID, err)
	}

	d.Set("addresses", sharedIPs)
	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_instance_shared_ips")

	linodeID := d.Get("linode_id").(int)
	d.SetId(strconv.Itoa(linodeID))

	return updateResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_instance_shared_ips")

	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	ctx = tflog.SetField(ctx, "linode_id", linodeID)

	ips := helper.ExpandStringSet(d.Get("addresses").(*schema.Set))

	if d.HasChange("addresses") {
		tflog.Info(ctx, "Updating shared IP addresses", map[string]any{
			"ips": ips,
		})

		err := client.ShareIPAddresses(ctx, linodego.IPAddressesShareOptions{
			LinodeID: linodeID,
			IPs:      ips,
		})
		if err != nil {
			return diag.Errorf("failed to update ips for linode %d: %s", linodeID, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_instance_shared_ips")

	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)

	err := client.ShareIPAddresses(ctx, linodego.IPAddressesShareOptions{
		LinodeID: linodeID,
		IPs:      []string{},
	})
	if err != nil {
		return diag.Errorf("failed to update ips for linode %d: %s", linodeID, err)
	}

	return nil
}

func GetSharedIPsForLinode(ctx context.Context, client linodego.Client, linodeID int) ([]string, error) {
	networking, err := client.GetInstanceIPAddresses(ctx, linodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance (%d) networking: %s", linodeID, err)
	}

	result := make([]string, 0)
	for _, ip := range networking.IPv4.Shared {
		result = append(result, ip.Address)
	}

	for _, ip := range networking.IPv6.Global {
		// BGP ips will not have a route target
		if ip.RouteTarget != "" {
			continue
		}

		result = append(result, ip.Range)
	}

	return result, nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": d.Get("linode_id").(int),
		"id":        d.Id(),
	})
}
