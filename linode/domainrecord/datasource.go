package domainrecord

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

	domainID := d.Get("domain_id").(int)
	recordName := d.Get("name").(string)
	recordID := d.Get("id").(int)

	if recordName == "" && recordID == 0 {
		return diag.Errorf("Record name or ID is required")
	}

	var record *linodego.DomainRecord

	if recordID != 0 {
		rec, err := client.GetDomainRecord(ctx, domainID, recordID)
		if err != nil {
			return diag.Errorf("Error fetching domain record: %v", err)
		}
		record = rec
	} else if recordName != "" {
		filter, _ := json.Marshal(map[string]interface{}{"name": recordName})
		records, err := client.ListDomainRecords(ctx, domainID, linodego.NewListOptions(0, string(filter)))
		if err != nil {
			return diag.Errorf("Error listing domain records: %v", err)
		}
		if len(records) > 0 {
			record = &records[0]
		}
	}

	if record != nil {
		d.SetId(strconv.Itoa(recordID))
		d.Set("id", record.ID)
		d.Set("name", record.Name)
		d.Set("type", record.Type)
		d.Set("ttl_sec", record.TTLSec)
		d.Set("target", record.Target)
		d.Set("priority", record.Priority)
		d.Set("protocol", record.Protocol)
		d.Set("weight", record.Weight)
		d.Set("port", record.Port)
		d.Set("service", record.Service)
		d.Set("tag", record.Tag)
		return nil
	}

	d.SetId("")

	return diag.Errorf(`Domain record "%s" for domain %d was not found`, recordName, domainID)
}
