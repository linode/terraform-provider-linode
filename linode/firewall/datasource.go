package firewall

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func dataSourceFirewallRules() *schema.Resource {
	return &schema.Resource{
		Schema: dataSourceRuleSchema,
	}
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id := d.Get("id").(int)

	firewall, err := client.GetFirewall(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get firewall %d: %s", id, err)
	}

	rules, err := client.GetFirewallRules(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get firewall rules %d: %s", id, err)
	}

	devices, err := client.ListFirewallDevices(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get firewall devices %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("label", firewall.Label)
	d.Set("tags", firewall.Tags)
	d.Set("disabled", firewall.Status == linodego.FirewallDisabled)
	d.Set("inbound", flattenFirewallRules(rules.Inbound))
	d.Set("inbound_policy", rules.InboundPolicy)
	d.Set("outbound", flattenFirewallRules(rules.Outbound))
	d.Set("outbound_policy", rules.OutboundPolicy)
	d.Set("status", firewall.Status)
	d.Set("linodes", flattenFirewallLinodes(devices))
	d.Set("devices", flattenFirewallDevices(devices))

	return nil
}
