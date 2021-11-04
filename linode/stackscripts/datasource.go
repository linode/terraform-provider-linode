package stackscripts

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/stackscript"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	latestFlag := d.Get("latest").(bool)

	filterID, err := helper.GetFilterID(d)
	if err != nil {
		return diag.Errorf("failed to generate filter id: %s", err)
	}

	filter, err := helper.ConstructFilterString(d, scriptValueToFilterType)
	if err != nil {
		return diag.Errorf("failed to construct filter: %s", err)
	}

	scripts, err := client.ListStackscripts(ctx, &linodego.ListOptions{
		Filter: filter,
	})
	if err != nil {
		return diag.Errorf("failed to list linode scripts: %s", err)
	}

	scriptsFlattened := make([]interface{}, len(scripts))
	for i, image := range scripts {
		scriptsFlattened[i] = flattenStackScript(&image)
	}

	scriptsFiltered, err := helper.FilterResults(d, scriptsFlattened)
	if err != nil {
		return diag.Errorf("failed to filter returned scripts: %s", err)
	}

	if latestFlag {
		latestScript := helper.GetLatestCreated(scriptsFiltered)

		if latestScript != nil {
			scriptsFiltered = []map[string]interface{}{latestScript}
		}
	}

	d.SetId(filterID)
	d.Set("stackscripts", scriptsFiltered)

	return nil
}

func flattenStackScript(script *linodego.Stackscript) map[string]interface{} {
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

	result["user_defined_fields"] = stackscript.GetStackScriptUserDefinedFields(script)

	return result
}

func scriptValueToFilterType(filterName, value string) (interface{}, error) {
	switch filterName {
	case "deprecated", "mine":
		return strconv.ParseBool(value)

	case "deployments_total":
		return strconv.Atoi(value)
	}

	return value, nil
}
