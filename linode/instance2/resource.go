package instance2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"log"
	"strconv"
	"time"
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
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance ID %s as int: %s", d.Id(), err)
	}

	inst, err := client.GetInstance(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Instance ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Instance: %s", err)
	}

	ips, err := client.GetInstanceIPAddresses(ctx, inst.ID)
	if err != nil {
		return diag.Errorf("failed to get ips for instance %d: %s", inst.ID, err)
	}

	d.Set("label", inst.Label)
	d.Set("region", inst.Region)
	d.Set("type", inst.Type)
	d.Set("tags", inst.Tags)
	d.Set("private_ip", len(ips.IPv4.Private) > 0)

	d.Set("ipv4_public", flattenInstanceIPs(ips.IPv4.Public))
	d.Set("ipv4_private", flattenInstanceIPs(ips.IPv4.Private))
	d.Set("ipv4_shared", flattenInstanceIPs(ips.IPv4.Shared))
	d.Set("ipv4_reserved", flattenInstanceIPs(ips.IPv4.Reserved))

	d.Set("ipv6_slaac", ips.IPv6.SLAAC.Address)
	d.Set("ipv6_link_local", ips.IPv6.LinkLocal.Address)
	d.Set("ipv6_global", flattenInstanceIPv6(ips.IPv6.Global))

	d.Set("specs", flattenInstanceSpecs(inst))

	d.Set("hypervisor", inst.Hypervisor)
	d.Set("status", inst.Status)
	d.Set("created", inst.Created.Format(time.RFC3339))
	d.Set("updated", inst.Updated.Format(time.RFC3339))
	d.Set("watchdog_enabled", inst.WatchdogEnabled)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	inst, err := client.CreateInstance(ctx, linodego.InstanceCreateOptions{
		Region:    d.Get("region").(string),
		Type:      d.Get("type").(string),
		Label:     d.Get("label").(string),
		Tags:      helper.ExpandStringSet(d.Get("tags").(*schema.Set)),
		PrivateIP: d.Get("private_ip").(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create linode instance: %s", err)
	}

	d.SetId(strconv.Itoa(inst.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance id %s as int: %s", d.Id(), err)
	}

	putRequest := linodego.InstanceUpdateOptions{}
	shouldUpdate := false

	// PUT updates
	if d.HasChange("tags") {
		newTags := helper.ExpandStringSet(d.Get("tags").(*schema.Set))
		putRequest.Tags = &newTags
		shouldUpdate = true
	}

	if d.HasChange("label") {
		putRequest.Label = d.Get("label").(string)
		shouldUpdate = true
	}

	if shouldUpdate {
		if _, err := client.UpdateInstance(ctx, int(id), putRequest); err != nil {
			return diag.Errorf("failed to update instance %d: %s", id, err)
		}
	}

	// Manual updates
	if d.HasChange("private_ip") {
		privateIP := d.Get("private_ip").(bool)
		if !privateIP {
			return diag.Errorf("private_ip cannot be disabled once enabled")
		}

		if _, err := client.AddInstanceIPAddress(ctx, int(id), false); err != nil {
			return diag.Errorf("failed to allocate private ip for instance %d: %s", id, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id64, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Instance id %s as int", d.Id())
	}

	err = client.DeleteInstance(ctx, int(id64))
	if err != nil {
		return diag.Errorf("Error deleting Linode Instance %d: %s", id64, err)
	}
	d.SetId("")
	return nil
}

func flattenInstanceIPs(ips []*linodego.InstanceIP) []string {
	result := make([]string, 0)

	for _, ip := range ips {
		if ip == nil {
			continue
		}

		result = append(result, ip.Address)
	}

	return result
}

func flattenInstanceIPv6(ips []linodego.IPv6Range) []string {
	result := make([]string, len(ips))

	for _, ip := range ips {
		result = append(result, ip.Range)
	}

	return result
}

func flattenInstanceSpecs(instance *linodego.Instance) []map[string]int {
	return []map[string]int{{
		"vcpus":    instance.Specs.VCPUs,
		"disk":     instance.Specs.Disk,
		"memory":   instance.Specs.Memory,
		"transfer": instance.Specs.Transfer,
	}}
}
