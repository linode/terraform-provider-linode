//go:build integration || maintenancepolicies

package maintenancepolicies_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/maintenancepolicies/tmpl"
)

func TestAccDataSourceMaintenancePolicies_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_maintenance_policies.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("slug"),
						knownvalue.StringExact("linode/migrate"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact("Migrate"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.StringExact("Migrates the Linode to a new host while it remains fully operational. Recommended for maximizing availability."),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.StringExact("migrate"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("notification_period_sec"),
						knownvalue.Int64Func(func(value int64) error {
							if value <= 0 {
								return fmt.Errorf("expected non-zero notification_period_sec")
							}

							return nil
						}),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("maintenance_policies").AtSliceIndex(0).AtMapKey("is_default"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}
