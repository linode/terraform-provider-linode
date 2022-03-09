package stackscripts

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/stackscript"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	results, err := filterConfig.FilterDataSource(ctx, d, meta, listStackScripts, flattenStackScript)
	if err != nil {
		return nil
	}

	results = filterConfig.FilterLatest(d, results)

	d.Set("stackscripts", results)

	return nil
}

func listStackScripts(
	ctx context.Context, d *schema.ResourceData, client *linodego.Client, options *linodego.ListOptions) ([]interface{}, error) {
	scripts, err := client.ListStackscripts(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(scripts))

	for i, v := range scripts {
		result[i] = v
	}

	return result, nil
}

func flattenStackScript(data interface{}) map[string]interface{} {
	script := data.(linodego.Stackscript)

	result := make(map[string]interface{})

	result["id"] = script.ID
	result["label"] = script.Label
	result["script"] = script.Script
	result["description"] = script.Description
	result["rev_note"] = script.RevNote
	result["is_public"] = script.IsPublic
	result["images"] = script.Images
	result["user_gravatar_id"] = script.UserGravatarID
	result["deployments_active"] = script.DeploymentsActive
	result["deployments_total"] = script.DeploymentsTotal
	result["username"] = script.Username

	if script.Created != nil {
		result["created"] = script.Created.Format(time.RFC3339)
	}

	if script.Updated != nil {
		result["updated"] = script.Updated.Format(time.RFC3339)
	}

	result["user_defined_fields"] = stackscript.GetStackScriptUserDefinedFields(&script)

	return result
}
