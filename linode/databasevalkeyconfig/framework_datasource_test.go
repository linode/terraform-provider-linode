//go:build integration || databasevalkeyconfig

package databasevalkeyconfig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/databasevalkeyconfig/tmpl"
)

func TestAccDataSourceDatabaseValkeyConfig_basic(t *testing.T) {
	t.Parallel()
	const resourceName = "data.linode_database_valkey_config.foobar"
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}
	for _, option := range []string{
		"backup_hour", "backup_minute", "frequent_snapshots", "valkey_acl_channels_default",
		"valkey_active_expire_effort", "valkey_activedefrag", "valkey_lfu_decay_time", "valkey_lfu_log_factor",
		"valkey_maxmemory_policy", "valkey_number_of_databases", "valkey_persistence",
		"valkey_pubsub_client_output_buffer_limit", "valkey_timeout",
	} {
		for _, field := range []string{"description", "requires_restart", "type"} {
			checks = append(checks, statecheck.ExpectKnownValue(resourceName,
				tfjsonpath.New(option).AtSliceIndex(0).AtMapKey(field), knownvalue.NotNull()))
		}
	}
	for _, option := range []string{"valkey_acl_channels_default", "valkey_maxmemory_policy", "valkey_persistence"} {
		checks = append(checks, statecheck.ExpectKnownValue(resourceName,
			tfjsonpath.New(option).AtSliceIndex(0).AtMapKey("enum"), knownvalue.NotNull()))
	}
	checks = append(checks,
		statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("backup_hour").AtSliceIndex(0).AtMapKey("type"),
			knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("integer"), knownvalue.StringExact("null")})),
		statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("frequent_snapshots").AtSliceIndex(0).AtMapKey("type"),
			knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("boolean")})),
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: tmpl.DataBasic(t), ConfigStateChecks: checks}},
	})
}
