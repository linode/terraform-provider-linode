//go:build integration || locks

package locks_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/locks/tmpl"
)

const (
	testLocksDataName = "data.linode_locks.test"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceLocks_basic(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel, testRegion),
				ConfigStateChecks: []statecheck.StateCheck{
					// TODO: Check ListSizeAtLeast 1 if ListSizeAtLeast is implemented on the upstream testing module
					// https://github.com/hashicorp/terraform-plugin-testing/issues/418
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("entity_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("entity_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("lock_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("entity_label"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("entity_url"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDataSourceLocks_filter(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	lockType := "cannot_delete"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilter(t, instanceLabel, testRegion, lockType),
				ConfigStateChecks: []statecheck.StateCheck{
					// TODO: Check ListSizeAtLeast 1 if ListSizeAtLeast is implemented on the upstream testing module
					// https://github.com/hashicorp/terraform-plugin-testing/issues/418
					statecheck.ExpectKnownValue(testLocksDataName, tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("entity_type"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						testLocksDataName,
						tfjsonpath.New("locks").AtSliceIndex(0).AtMapKey("lock_type"),
						knownvalue.StringExact(lockType),
					),
				},
			},
		},
	})
}
