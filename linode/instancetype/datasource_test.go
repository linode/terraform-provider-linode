//go:build integration || instancetype

package instancetype_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetype/tmpl"
)

func TestAccDataSourceLinodeInstanceType_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_type.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fatal(err)
	}

	// Resolve a type with region-specific pricing
	allTypes, err := client.ListTypes(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to list regions: %s", err)
	}

	var targetType linodego.LinodeType
	for _, v := range allTypes {
		if len(v.RegionPrices) > 0 && v.RegionPrices[0].Hourly > 0 {
			targetType = v
			break
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, targetType.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", targetType.ID),
					resource.TestCheckResourceAttr(resourceName, "label", targetType.Label),
					resource.TestCheckResourceAttr(
						resourceName,
						"disk",
						strconv.FormatInt(int64(targetType.Disk), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"class",
						string(targetType.Class),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"memory",
						strconv.FormatInt(int64(targetType.Memory), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"vcpus",
						strconv.FormatInt(int64(targetType.VCPUs), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"network_out",
						strconv.FormatInt(int64(targetType.NetworkOut), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"price.0.hourly",
						strconv.FormatFloat(float64(targetType.Price.Hourly), 'f', -1, 64)),
					resource.TestCheckResourceAttr(
						resourceName,
						"price.0.monthly",
						strconv.FormatFloat(float64(targetType.Price.Monthly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.price.0.hourly",
						strconv.FormatFloat(float64(targetType.Addons.Backups.Price.Hourly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.price.0.monthly",
						strconv.FormatFloat(float64(targetType.Addons.Backups.Price.Monthly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"region_prices.0.monthly",
						strconv.FormatFloat(float64(targetType.RegionPrices[0].Monthly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"region_prices.0.hourly",
						strconv.FormatFloat(float64(targetType.RegionPrices[0].Hourly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.region_prices.0.monthly",
						strconv.FormatFloat(float64(targetType.Addons.Backups.RegionPrices[0].Monthly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.region_prices.0.hourly",
						strconv.FormatFloat(float64(targetType.Addons.Backups.RegionPrices[0].Hourly), 'f', -1, 64),
					),
				),
			},
		},
	})
}
