package databases

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listDatabases, flattenDatabase)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("databases", results)

	return nil
}

func listDatabases(
	ctx context.Context, d *schema.ResourceData,
	client *linodego.Client, options *linodego.ListOptions) ([]interface{}, error) {
	dbs, err := client.ListDatabases(ctx, options)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(dbs))

	for i, v := range dbs {
		result[i] = v
	}

	return result, nil
}

func flattenDatabase(data interface{}) map[string]interface{} {
	db := data.(linodego.Database)

	result := make(map[string]interface{})

	result["id"] = db.ID
	result["status"] = db.Status
	result["label"] = db.Label
	result["host_primary"] = db.Hosts.Primary
	result["host_secondary"] = db.Hosts.Secondary
	result["region"] = db.Region
	result["type"] = db.Type
	result["engine"] = db.Engine
	result["version"] = db.Version
	result["cluster_size"] = db.ClusterSize
	result["replication_type"] = db.ReplicationType
	result["ssl_connection"] = db.SSLConnection
	result["encrypted"] = db.Encrypted
	result["allow_list"] = db.AllowList
	result["instance_uri"] = db.InstanceURI
	result["created"] = db.Created.Format(time.RFC3339)
	result["updated"] = db.Updated.Format(time.RFC3339)

	return result
}
