package instancesharedips_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancesharedips"
	"github.com/linode/terraform-provider-linode/linode/instancesharedips/tmpl"
)

const (
	resourcePrimaryNode    = "linode_instance.primary"
	resourceSecondaryNode  = "linode_instance.secondary"
	resourcePrimaryShare   = "linode_instance_shared_ips.share-primary"
	resourceSecondaryShare = "linode_instance_shared_ips.share-secondary"
)

func TestAccInstanceSharedIPs_update(t *testing.T) {
	t.Parallel()

	var primaryInstance, secondaryInstance linodego.Instance

	name := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.SingleNode(t, name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resourcePrimaryNode, &primaryInstance),
					acceptance.CheckInstanceExists(resourceSecondaryNode, &secondaryInstance),

					checkInstanceSharedIPCount(resourcePrimaryNode, 1),
					checkInstanceSharedIPCount(resourceSecondaryNode, 0),

					resource.TestCheckResourceAttr(resourcePrimaryShare, "addresses.#", "1"),
				),
			},
			{
				Config: tmpl.DualNode(t, name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resourcePrimaryNode, &primaryInstance),
					acceptance.CheckInstanceExists(resourceSecondaryNode, &secondaryInstance),

					checkInstanceSharedIPCount(resourcePrimaryNode, 1),
					checkInstanceSharedIPCount(resourceSecondaryNode, 1),

					resource.TestCheckResourceAttr(resourcePrimaryShare, "addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceSecondaryShare, "addresses.#", "1"),
				),
			},
			{
				Config: tmpl.SingleNode(t, name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resourcePrimaryNode, &primaryInstance),
					acceptance.CheckInstanceExists(resourceSecondaryNode, &secondaryInstance),

					checkInstanceSharedIPCount(resourcePrimaryNode, 1),
					checkInstanceSharedIPCount(resourceSecondaryNode, 0),

					resource.TestCheckResourceAttr(resourcePrimaryShare, "addresses.#", "1"),
				),
			},
		},
	})
}

func checkInstanceSharedIPCount(name string, length int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		ips, err := instancesharedips.GetSharedIPsForLinode(context.Background(), client, id)
		if err != nil {
			return err
		}

		if len(ips) != length {
			return fmt.Errorf("lengths do not match: %d != %d", len(ips), length)
		}

		return nil
	}
}
