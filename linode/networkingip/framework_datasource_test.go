//go:build integration || networkingip

package networkingip_test

import (
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"linodes"}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccDataSourceNetworkingIP_basic(t *testing.T) {
	t.Parallel()

	resourceName := "linode_instance.foobar"
	dataResourceName := "data.linode_networking_ip.foobar"

	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion),
			},
			{
				Config: tmpl.DataBasic(t, label, testRegion),
				Check: resource.ComposeTestCheckFunc(
					// statechecks can't compare int linode_id with string id without implementing a custom comparer.
					// Keep this legacy check for now.
					resource.TestCheckResourceAttrPair(dataResourceName, "linode_id", resourceName, "id"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.CompareValuePairs(dataResourceName, tfjsonpath.New("address"), resourceName, tfjsonpath.New("ipv4").AtSliceIndex(0), compare.ValuesSame()),
					statecheck.CompareValuePairs(dataResourceName, tfjsonpath.New("region"), resourceName, tfjsonpath.New("region"), compare.ValuesSame()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("gateway"), knownvalue.StringRegexp(regexp.MustCompile(`\.1$`))),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("type"), knownvalue.StringExact("ipv4")),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("public"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("reserved"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("prefix"), knownvalue.Int64Exact(24)),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("rdns"), knownvalue.StringRegexp(regexp.MustCompile(`.ip.linodeusercontent.com$`))),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("vpc_nat_1_1"), knownvalue.Null()),
					statecheck.ExpectKnownValue(dataResourceName, tfjsonpath.New("interface_id"), knownvalue.NotNull()),
				},
			},
		},
	})
}
