//go:build unit

package databasevalkeyconfig

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
)

func TestModel_ParseValkeyConfig(t *testing.T) {
	var config linodego.ValkeyDatabaseConfigInfo
	require.NoError(t, json.Unmarshal([]byte(`{
		"backup_hour":{"description":"Backup hour","example":0,"minimum":0,"maximum":23,"type":["integer","null"]},
		"backup_minute":{"example":30,"minimum":0,"maximum":59,"type":["integer","null"]},
		"frequent_snapshots":{"default":false,"type":"boolean"},
		"valkey_acl_channels_default":{"enum":["allchannels","resetchannels"],"type":"string"},
		"valkey_active_expire_effort":{"default":1,"minimum":1,"maximum":10,"type":"integer"},
		"valkey_activedefrag":{"default":true,"type":"boolean"},
		"valkey_lfu_decay_time":{"default":1,"minimum":1,"maximum":120,"type":"integer"},
		"valkey_lfu_log_factor":{"default":10,"minimum":0,"maximum":100,"type":"integer"},
		"valkey_maxmemory_policy":{"default":"noeviction","enum":["noeviction","allkeys-lru"],"type":["string","null"]},
		"valkey_number_of_databases":{"example":16,"minimum":1,"maximum":128,"requires_restart":true,"type":"integer"},
		"valkey_persistence":{"enum":["off","rdb"],"type":"string"},
		"valkey_pubsub_client_output_buffer_limit":{"example":32,"minimum":32,"maximum":262144,"type":"integer"},
		"valkey_timeout":{"default":0,"minimum":0,"maximum":2073600,"type":"integer"}
	}`), &config))
	var model DataSourceModel
	var diags diag.Diagnostics
	model.ParseValkeyConfig(t.Context(), &config, &diags)
	require.False(t, diags.HasError(), "%v", diags)

	state := tfsdk.State{Schema: frameworkDataSourceSchema}
	diags = state.Set(t.Context(), &model)
	require.False(t, diags.HasError(), "%v", diags)
	var value map[string]tftypes.Value
	require.NoError(t, state.Raw.As(&value))
	require.Len(t, value, 14)
	for name := range frameworkDataSourceSchema.Attributes {
		if name == "id" {
			continue
		}
		require.False(t, value[name].IsNull(), name)
	}
	backup := model.BackupHour.Elements()[0].(types.Object).Attributes()
	require.Equal(t, types.Int64Value(0), backup["example"])
	require.Equal(t, types.Int64Value(23), backup["maximum"])
	require.NotContains(t, backup, "default")
	require.Equal(t, []attr.Value{types.StringValue("integer"), types.StringValue("null")}, backup["type"].(types.List).Elements())
	frequent := model.FrequentSnapshots.Elements()[0].(types.Object).Attributes()
	require.Equal(t, types.BoolValue(false), frequent["default"])
	require.NotContains(t, frequent, "minimum")
	policy := model.ValkeyMaxmemoryPolicy.Elements()[0].(types.Object).Attributes()
	require.Equal(t, types.StringValue("noeviction"), policy["default"])
	require.Len(t, policy["enum"].(types.List).Elements(), 2)
	require.Equal(t, types.BoolValue(true), model.ValkeyNumberOfDatabases.Elements()[0].(types.Object).Attributes()["requires_restart"])
	require.Equal(t, types.Int64Value(0), model.ValkeyTimeout.Elements()[0].(types.Object).Attributes()["default"])

	var copied DataSourceModel
	copied.ParseValkeyConfig(t.Context(), &config, &diags)
	require.Equal(t, model, copied)
	config.BackupHour.Example = 1
	copied.ParseValkeyConfig(t.Context(), &config, &diags)
	require.NotEqual(t, model.ID, copied.ID)
	require.False(t, diags.HasError(), "%v", diags)
}
