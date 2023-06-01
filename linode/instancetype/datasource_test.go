package instancetype_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/instancetype/tmpl"
)

func TestAccDataSourceLinodeInstanceType_basic(t *testing.T) {
	t.Parallel()

	instanceTypeID := "g6-standard-2"
	resourceName := "data.linode_instance_type.foobar"

	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		t.Fatal(err)
	}

	typeInfo, err := client.GetType(context.Background(), instanceTypeID)
	if err != nil {
		t.Fatalf("failed to get instance type %s: %s", instanceTypeID, err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceTypeID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", instanceTypeID),
					resource.TestCheckResourceAttr(resourceName, "label", typeInfo.Label),
					resource.TestCheckResourceAttr(
						resourceName,
						"disk",
						strconv.FormatInt(int64(typeInfo.Disk), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"class",
						string(typeInfo.Class),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"memory",
						strconv.FormatInt(int64(typeInfo.Memory), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"vcpus",
						strconv.FormatInt(int64(typeInfo.VCPUs), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"network_out",
						strconv.FormatInt(int64(typeInfo.NetworkOut), 10),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"price.0.hourly",
						strconv.FormatFloat(float64(typeInfo.Price.Hourly), 'f', -1, 64)),
					resource.TestCheckResourceAttr(
						resourceName,
						"price.0.monthly",
						strconv.FormatFloat(float64(typeInfo.Price.Monthly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.price.0.hourly",
						strconv.FormatFloat(float64(typeInfo.Addons.Backups.Price.Hourly), 'f', -1, 64),
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"addons.0.backups.0.price.0.monthly",
						strconv.FormatFloat(float64(typeInfo.Addons.Backups.Price.Monthly), 'f', -1, 64),
					),
				),
			},
		},
	})
}
