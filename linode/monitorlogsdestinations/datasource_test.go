//go:build integration || monitorlogsdestinations

package monitorlogsdestinations_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsdestinations/tmpl"
)

func TestAccDataSourceLogsDestinations_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityObjectStorage}, "core")
	if err != nil {
		t.Skipf("could not get region with Object Storage: %s", err)
	}

	dataName := "data.linode_monitor_logs_destinations.foobar"

	acceptance.RunTestWithRetries(t, 5, func(t *acceptance.WrappedT) {
		label := acctest.RandomWithPrefix("tf-test")

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: tmpl.DataBasic(t, label, region),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(dataName,
							tfjsonpath.New("destinations"),
							knownvalue.ListSizeExact(1),
						),
						statecheck.ExpectKnownValue(dataName,
							tfjsonpath.New("destinations").AtSliceIndex(0).AtMapKey("label"),
							knownvalue.StringExact(label),
						),
						statecheck.ExpectKnownValue(dataName,
							tfjsonpath.New("destinations").AtSliceIndex(0).AtMapKey("type"),
							knownvalue.StringExact("akamai_object_storage"),
						),
						statecheck.ExpectKnownValue(dataName,
							tfjsonpath.New("destinations").AtSliceIndex(0).AtMapKey("status"),
							knownvalue.NotNull(),
						),
					},
				},
			},
		})
	})
}
