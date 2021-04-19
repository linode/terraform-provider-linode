package linode

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"time"

	"fmt"
	"testing"
)

func TestAccDataSourceLinodeVLANs_basic(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	vlanName := acctest.RandomWithPrefix("tf-test")
	resourceName := "data.linode_vlans.foolan"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeVLANsBasic(instanceName, vlanName),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*ProviderMeta).Client
					if _, err := waitForVLANWithLabel(client, vlanName, 30); err != nil {
						t.Fatal(err)
					}
				},
				Config: testDataSourceLinodeVLANsBasic(instanceName, vlanName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vlans.0.label", vlanName),
					resource.TestCheckResourceAttr(resourceName, "vlans.0.region", "us-southeast"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.linodes.#"),
				),
			},
		},
	})
}

// waitForVLANWithLabel polls for a VLAN with the given label to exist
// This is necessary in this context because it is not guaranteed that a VLAN will
// be created immediately after the instance is created.
func waitForVLANWithLabel(client linodego.Client, label string, timeoutSeconds int) (*linodego.VLAN, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSeconds))
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			vlans, err := client.ListVLANs(ctx,
				&linodego.ListOptions{Filter: fmt.Sprintf("{\"label\": \"%s\"}", label)})
			if err != nil {
				return nil, err
			}

			if len(vlans) > 0 {
				return &vlans[0], err
			}

		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for VLAN %s: %s", label, ctx.Err())
		}
	}
}

func testDataSourceLinodeVLANsBasic(instanceName, vlanName string) string {
	return fmt.Sprintf(`
resource "linode_instance" "fooinst" {
	label = "%s"
	type = "g6-standard-1"
	image = "linode/alpine3.13"
	region = "us-southeast"

	interface {
		label = "%s"
		purpose = "vlan"
	}
}

data "linode_vlans" "foolan" {
	filter {
		name = "label"
		values = ["%s"]
	}
}`, instanceName, vlanName, vlanName)
}
