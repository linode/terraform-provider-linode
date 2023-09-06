package firewall

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func resourceFirewallRules() *schema.Resource {
	return &schema.Resource{
		Schema: resourceRuleSchema,
	}
}

func resourceFirewallDevice() *schema.Resource {
	return &schema.Resource{
		Schema: resourceDeviceSchema,
	}
}

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

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse Firewall %s as int: %s", d.Id(), err)
	}
	firewall, err := client.GetFirewall(ctx, id)
	if err != nil {
		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code == 404 {
			log.Printf("[WARN] removing Linode Firewall ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get firewall %d: %s", id, err)
	}

	rules, err := client.GetFirewallRules(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get rules for firewall %d: %s", id, err)
	}

	devices, err := client.ListFirewallDevices(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get devices for firewall %d: %s", id, err)
	}

	d.Set("label", firewall.Label)
	d.Set("disabled", firewall.Status == linodego.FirewallDisabled)
	d.Set("tags", firewall.Tags)
	d.Set("status", firewall.Status)
	d.Set("created", firewall.Created.Format(helper.TIME_FORMAT))
	d.Set("updated", firewall.Updated.Format(helper.TIME_FORMAT))
	d.Set("inbound", flattenFirewallRules(rules.Inbound))
	d.Set("outbound", flattenFirewallRules(rules.Outbound))
	d.Set("inbound_policy", firewall.Rules.InboundPolicy)
	d.Set("outbound_policy", firewall.Rules.OutboundPolicy)
	d.Set("linodes", flattenFirewallDeviceIDs(devices, linodego.FirewallDeviceLinode))
	d.Set("nodebalancers", flattenFirewallDeviceIDs(devices, linodego.FirewallDeviceNodeBalancer))
	d.Set("devices", flattenFirewallDevices(devices))
	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.FirewallCreateOptions{
		Label: d.Get("label").(string),
		Tags:  helper.ExpandStringSet(d.Get("tags").(*schema.Set)),
	}

	createOpts.Devices.Linodes = helper.ExpandIntSet(d.Get("linodes").(*schema.Set))
	createOpts.Devices.NodeBalancers = helper.ExpandIntSet(d.Get("nodebalancers").(*schema.Set))
	createOpts.Rules.Inbound = expandFirewallRules(d.Get("inbound").([]any))
	createOpts.Rules.InboundPolicy = d.Get("inbound_policy").(string)
	createOpts.Rules.Outbound = expandFirewallRules(d.Get("outbound").([]any))
	createOpts.Rules.OutboundPolicy = d.Get("outbound_policy").(string)

	firewall, err := client.CreateFirewall(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create Firewall: %s", err)
	}
	d.SetId(strconv.Itoa(firewall.ID))

	if d.Get("disabled").(bool) {
		if _, err := client.UpdateFirewall(ctx, firewall.ID, linodego.FirewallUpdateOptions{
			Status: linodego.FirewallDisabled,
		}); err != nil {
			return diag.Errorf("failed to disable firewall %d: %s", firewall.ID, err)
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse Firewall %s as int: %s", d.Id(), err)
	}

	if d.HasChanges("label", "tags", "disabled") {
		updateOpts := linodego.FirewallUpdateOptions{}
		if d.HasChange("label") {
			updateOpts.Label = d.Get("label").(string)
		}
		if d.HasChange("tags") {
			tags := helper.ExpandStringSet(d.Get("tags").(*schema.Set))
			updateOpts.Tags = &tags
		}
		if d.HasChange("disabled") {
			updateOpts.Status = expandFirewallStatus(d.Get("disabled"))
		}

		if _, err := client.UpdateFirewall(ctx, id, updateOpts); err != nil {
			return diag.Errorf("failed to update firewall %d: %s", id, err)
		}
	}

	inboundRules := expandFirewallRules(d.Get("inbound").([]any))
	outboundRules := expandFirewallRules(d.Get("outbound").([]any))
	ruleSet := linodego.FirewallRuleSet{
		Inbound:        inboundRules,
		InboundPolicy:  d.Get("inbound_policy").(string),
		Outbound:       outboundRules,
		OutboundPolicy: d.Get("outbound_policy").(string),
	}
	if _, err := client.UpdateFirewallRules(ctx, id, ruleSet); err != nil {
		return diag.Errorf("failed to update rules for firewall %d: %s", id, err)
	}

	linodes, linodesOk := d.GetOk("linodes")
	nodebalancers, nodebalancersOk := d.GetOk("nodebalancers")

	if linodesOk || nodebalancersOk {
		assignments := make([]firewallDeviceAssignment, 0)

		for _, entityID := range helper.ExpandIntSet(linodes.(*schema.Set)) {
			assignments = append(assignments, firewallDeviceAssignment{
				ID:   entityID,
				Type: linodego.FirewallDeviceLinode,
			})
		}

		for _, entityID := range helper.ExpandIntSet(nodebalancers.(*schema.Set)) {
			assignments = append(assignments, firewallDeviceAssignment{
				ID:   entityID,
				Type: linodego.FirewallDeviceNodeBalancer,
			})
		}

		if err := updateFirewallDevices(ctx, d, client, id, assignments); err != nil {
			return diag.Errorf("failed to update firewall devices: %s", err)
		}
	}

	return nil
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed to parse Firewall %s as int: %s", d.Id(), err)
	}

	if err := client.DeleteFirewall(ctx, id); err != nil {
		return diag.Errorf("failed to delete Firewall %d: %s", id, err)
	}
	return nil
}
