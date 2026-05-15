//go:build integration || monitorlogsdestination

package monitorlogsdestination_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/monitorlogsdestination/tmpl"
)

func TestAccDataSourceLogsDestination_basic(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityObjectStorage}, "core")
	if err != nil {
		t.Skipf("could not get region with Object Storage: %s", err)
	}

	dataName := "data.linode_monitor_logs_destination.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("label"), knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("type"), knownvalue.StringExact("akamai_object_storage")),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("status"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("updated"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName, tfjsonpath.New("created_by"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataName,
						tfjsonpath.New("details").AtMapKey("bucket_name"),
						knownvalue.StringExact(label+"-bucket"),
					),
				},
			},
			{
				Config:      tmpl.DataNotFound(t),
				ExpectError: regexp.MustCompile(`\[404\]`),
			},
		},
	})
}
