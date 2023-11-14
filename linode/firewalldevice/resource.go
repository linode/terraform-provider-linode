package firewalldevice

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid firewall device ID: %v", err)
		}

		firewallID, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid firewall ID: %v", err)
		}

		d.SetId(s[1])
		d.Set("firewall_id", firewallID)
	}

	err := readResource(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import %v as firewall_device: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Firewall ID %s as int: %s", d.Id(), err)
	}

	firewallID := d.Get("firewall_id").(int)

	device, err := client.GetFirewallDevice(ctx, firewallID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Firewall Device ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Firewall Device: %s", err)
	}

	d.Set("entity_id", device.Entity.ID)
	d.Set("entity_type", device.Entity.Type)
	d.Set("created", device.Created.String())
	d.Set("updated", device.Updated.String())

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	firewallID := d.Get("firewall_id").(int)
	entityID := d.Get("entity_id").(int)
	entityType := d.Get("entity_type").(string)

	createOpts := linodego.FirewallDeviceCreateOptions{
		ID:   entityID,
		Type: linodego.FirewallDeviceType(entityType),
	}

	device, err := client.CreateFirewallDevice(ctx, firewallID, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Firewall Device: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", device.ID))

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Firewall Device ID %s as int: %s", d.Id(), err)
	}

	firewallID, ok := d.Get("firewall_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode Firewall Device ID %v as int", d.Get("firewall_id"))
	}

	err = client.DeleteFirewallDevice(ctx, firewallID, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode Firewall Device %d: %s", id, err)
	}
	return nil
}
