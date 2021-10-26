package vlan_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/vlan/tmpl"

	"context"
	"fmt"
	"testing"
	"time"
)

func TestAccDataSourceVLANs_basic(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	vlanName := acctest.RandomWithPrefix("tf-test")
	resourceName := "data.linode_vlans.foolan"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceName, vlanName),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
					if _, err := waitForVLANWithLabel(client, vlanName, 30); err != nil {
						t.Fatal(err)
					}
				},
				Config: tmpl.DataBasic(t, instanceName, vlanName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vlans.0.label", vlanName),
					resource.TestCheckResourceAttr(resourceName, "vlans.0.region", "us-southeast"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.created"),
					resource.TestCheckResourceAttrSet(resourceName, "vlans.0.linodes.#"),

					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order_by\":\"region\""),
					acceptance.CheckResourceAttrContains(resourceName, "id", "\"+order\":\"desc\""),
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
