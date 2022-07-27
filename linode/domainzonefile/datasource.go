package domainzonefile

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
	retryDuration := time.Duration(meta.(*helper.ProviderMeta).Config.EventPollMilliseconds)
	zf, err := getZoneFileRetry(ctx, &client, domainID, retryDuration)
	if err != nil {
		return diag.Errorf("%s", err)
	}

	d.SetId(strconv.Itoa(domainID))
	d.Set("domain_id", domainID)
	d.Set("zone_file", zf.ZoneFile)

	return nil
}

func getZoneFileRetry(ctx context.Context, client *linodego.Client,
	domainID int, retryDuration time.Duration,
) (*linodego.DomainZoneFile, error) {
	ticker := time.NewTicker(retryDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			zf, err := client.GetDomainZoneFile(ctx, domainID)
			if err != nil {
				return nil, fmt.Errorf("error fetching domain record: %v", err)
			}
			if len(zf.ZoneFile) > 0 {
				return zf, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("unable to fetch domain record")
		}
	}
}
