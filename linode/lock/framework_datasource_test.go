//go:build integration || lock

package lock_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/lock/tmpl"
)

const (
	testLockDataName = "data.linode_lock.test"
)

func TestAccDataSourceLock_basic(t *testing.T) {
	t.Parallel()

	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityLinodes}, "core")
	if err != nil {
		t.Fatal(err)
	}

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testLockDataName, "id"),
					resource.TestCheckResourceAttrSet(testLockDataName, "entity_id"),
					resource.TestCheckResourceAttr(testLockDataName, "entity_type", "linode"),
					resource.TestCheckResourceAttr(testLockDataName, "lock_type", "cannot_delete"),
					resource.TestCheckResourceAttrSet(testLockDataName, "entity_label"),
					resource.TestCheckResourceAttrSet(testLockDataName, "entity_url"),
				),
			},
		},
	})
}
