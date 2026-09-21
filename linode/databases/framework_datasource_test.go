//go:build integration || databases || dbaas_tests

package databases_test

import (
	"context"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/databases/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/helper/databaseshared"
)

var (
	testRegion    string
	engineVersion string
)

func init() {
	client, err := acceptance.GetTestClient()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := databaseshared.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failed to get db engine version: %s", err)
	}

	engineVersion = v.ID

	region, err := acceptance.GetRandomRegionWithCaps([]linodego.RegionCapability{linodego.CapabilityDBAAS}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceDatabases_byAttr(t *testing.T) {
	acceptance.LongRunningTest(t)
	t.Parallel()

	resourceName := "data.linode_databases.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ByLabel(t, engineVersion, dbName, dbName, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(dbName)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("cluster_size"),
						knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("encrypted"),
						knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("engine"),
						knownvalue.StringExact("mysql")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("platform"),
						knownvalue.StringExact("rdbms-default")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("region"),
						knownvalue.StringExact(testRegion)),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("type"),
						knownvalue.StringExact("g6-nanode-1")),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("allow_list"),
						knownvalue.SetSizeExact(0)),

					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("created"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("host_primary"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("status"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("updated"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("version"),
						knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("fork_restore_time"),
						knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("fork_source"),
						knownvalue.Null()),
				},
			},
			{
				Config: tmpl.ByLabel(t, engineVersion, dbName, "not"+dbName, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases"),
						knownvalue.ListSizeExact(0)),
				},
			},
			{
				Config: tmpl.ByEngine(t, engineVersion, dbName, "mysql", testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					// ListSizeAtLeast is not available upstream yet; index 0 implies size >= 1.
					// https://github.com/hashicorp/terraform-plugin-testing/issues/418
					statecheck.ExpectKnownValue(resourceName,
						tfjsonpath.New("databases").AtSliceIndex(0).AtMapKey("engine"),
						knownvalue.StringExact("mysql")),
				},
			},
		},
	})
}
