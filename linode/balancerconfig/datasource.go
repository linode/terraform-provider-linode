package balancerconfig

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id := d.Get("id").(int)
	nodebalancerID := d.Get("nodebalancer_id").(int)

	config, err := client.GetNodeBalancerConfig(ctx, nodebalancerID, id)
	if err != nil {
		return diag.Errorf("failed to get nodebalancer config %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(config.ID))
	d.Set("nodebalancer_id", config.NodeBalancerID)
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
