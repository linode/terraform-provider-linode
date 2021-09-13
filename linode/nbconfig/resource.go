package nbconfig

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

func resourceStatus() *schema.Resource {
	return &schema.Resource{
		Schema: resourceSchemaStatus,
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
			StateContext: importResource,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    ResourceNodeBalancerConfigV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceNodeBalancerConfigV0Upgrade,
				Version: 0,
			},
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		// Validate that this is an ID by making sure it can be converted into an int
		_, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid nodebalancer_config ID: %v", err)
		}

		nodebalancerID, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid nodebalancer ID: %v", err)
		}

		d.SetId(s[1])
		d.Set("nodebalancer_id", nodebalancerID)
	}

	err := readResource(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import %v as nodebalancer_config: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}

	config, err := client.GetNodeBalancerConfig(ctx, int(nodebalancerID), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing NodeBalancer Config ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode NodeBalancerConfig: %s", err)
	}

	d.Set("algorithm", config.Algorithm)
	d.Set("stickiness", config.Stickiness)
	d.Set("check", config.Check)
	d.Set("check_attempts", config.CheckAttempts)
	d.Set("check_body", config.CheckBody)
	d.Set("check_interval", config.CheckInterval)
	d.Set("check_timeout", config.CheckTimeout)
	d.Set("check_passive", config.CheckPassive)
	d.Set("check_path", config.CheckPath)
	d.Set("cipher_suite", config.CipherSuite)
	d.Set("port", config.Port)
	d.Set("protocol", config.Protocol)
	d.Set("proxy_protocol", config.ProxyProtocol)
	d.Set("ssl_fingerprint", config.SSLFingerprint)
	d.Set("ssl_commonname", config.SSLCommonName)
	d.Set("node_status", []map[string]interface{}{{
		"up":   config.NodesStatus.Up,
		"down": config.NodesStatus.Down,
	}})

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	nodebalancerID := d.Get("nodebalancer_id").(int)

	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Algorithm:     linodego.ConfigAlgorithm(d.Get("algorithm").(string)),
		Check:         linodego.ConfigCheck(d.Get("check").(string)),
		Stickiness:    linodego.ConfigStickiness(d.Get("stickiness").(string)),
		CheckAttempts: d.Get("check_attempts").(int),
		CheckBody:     d.Get("check_body").(string),
		CheckInterval: d.Get("check_interval").(int),
		CheckPath:     d.Get("check_path").(string),
		CheckTimeout:  d.Get("check_timeout").(int),
		Port:          d.Get("port").(int),
		Protocol:      linodego.ConfigProtocol(strings.ToLower(d.Get("protocol").(string))),
		ProxyProtocol: linodego.ConfigProxyProtocol(d.Get("proxy_protocol").(string)),
		SSLCert:       d.Get("ssl_cert").(string),
		SSLKey:        d.Get("ssl_key").(string),
	}

	if checkPassiveRaw, ok := d.GetOkExists("check_passive"); ok {
		checkPassive := checkPassiveRaw.(bool)
		createOpts.CheckPassive = &checkPassive
	}

	config, err := client.CreateNodeBalancerConfig(ctx, nodebalancerID, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode NodeBalancerConfig: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", config.ID))
	d.Set("nodebalancer_id", nodebalancerID)

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %s as int", d.Get("nodebalancer_id"))
	}

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		Algorithm:     linodego.ConfigAlgorithm(d.Get("algorithm").(string)),
		Check:         linodego.ConfigCheck(d.Get("check").(string)),
		Stickiness:    linodego.ConfigStickiness(d.Get("stickiness").(string)),
		CheckAttempts: d.Get("check_attempts").(int),
		CheckBody:     d.Get("check_body").(string),
		CheckInterval: d.Get("check_interval").(int),
		CheckPath:     d.Get("check_path").(string),
		CheckTimeout:  d.Get("check_timeout").(int),
		Port:          d.Get("port").(int),
		Protocol:      linodego.ConfigProtocol(strings.ToLower(d.Get("protocol").(string))),
		ProxyProtocol: linodego.ConfigProxyProtocol(d.Get("proxy_protocol").(string)),
		SSLCert:       d.Get("ssl_cert").(string),
		SSLKey:        d.Get("ssl_key").(string),
	}

	if ok := d.HasChange("check_passive"); ok {
		checkPassive := d.Get("check_passive").(bool)
		updateOpts.CheckPassive = &checkPassive
	}

	if _, err = client.UpdateNodeBalancerConfig(ctx, int(nodebalancerID), int(id), updateOpts); err != nil {
		return diag.Errorf("Error updating Nodebalancer %d Config %d: %s", int(nodebalancerID), int(id), err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	err = client.DeleteNodeBalancerConfig(ctx, nodebalancerID, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode NodeBalancerConfig %d: %s", id, err)
	}
	return nil
}

func ResourceNodeBalancerConfigV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_status": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func ResourceNodeBalancerConfigV0Upgrade(ctx context.Context,
	rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldTransfer, ok := rawState["node_status"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to upgrade state: node_status key does not exist")
	}

	newTransfer := []map[string]interface{}{
		{
			"down": 0,
			"up":   0,
		},
	}

	for key, val := range oldTransfer {
		val := val.(string)

		// This is necessary because it is possible old versions of the state have empty transfer fields
		// that must default to zero.
		if val == "" {
			continue
		}

		result, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("failed to parse state: %v", err)
		}

		newTransfer[0][key] = result
	}

	rawState["node_status"] = newTransfer
	return rawState, nil
}
