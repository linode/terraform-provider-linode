package domain

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      dataSourceSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	reqIDString := d.Get("id").(string)
	reqDomain := d.Get("domain").(string)

	if reqDomain == "" && reqIDString == "" {
		return diag.Errorf("Domain or Domain ID is required")
	}

	var domain *linodego.Domain

	d.SetId("")

	if reqIDString != "" {
		reqID, err := strconv.Atoi(reqIDString)
		if err != nil {
			return diag.Errorf("Domain ID %q must be numeric", reqIDString)
		}

		domain, err = client.GetDomain(ctx, reqID)
		if err != nil {
			return diag.Errorf("Error listing domain: %s", err)
		}
		if reqDomain != "" && domain.Domain != reqDomain {
			return diag.Errorf("Domain ID was found but did not match the requested Domain name")
		}
	} else if reqDomain != "" {
		filter, _ := json.Marshal(map[string]interface{}{"domain": reqDomain})
		domains, err := client.ListDomains(ctx, linodego.NewListOptions(0, string(filter)))
		if err != nil {
			return diag.Errorf("Error listing Domains: %s", err)
		}
		if len(domains) != 1 || domains[0].Domain != reqDomain {
			return diag.Errorf("Domain %s was not found", reqDomain)
		}
		domain = &domains[0]
	}

	if domain != nil {
		d.SetId(strconv.Itoa(domain.ID))
		d.Set("domain", domain.Domain)
		d.Set("type", domain.Type)
		d.Set("group", domain.Group)
		d.Set("status", domain.Status)
		d.Set("description", domain.Description)
		d.Set("master_ips", domain.MasterIPs)
		d.Set("axfr_ips", domain.AXfrIPs)
		d.Set("ttl_sec", domain.TTLSec)
		d.Set("retry_sec", domain.RetrySec)
		d.Set("expire_sec", domain.ExpireSec)
		d.Set("refresh_sec", domain.RefreshSec)
		d.Set("soa_email", domain.SOAEmail)
		d.Set("tags", domain.Tags)
		return nil
	}

	return diag.Errorf("Domain %s%s was not found", reqIDString, reqDomain)
}
