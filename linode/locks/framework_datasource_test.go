//go:build integration || locks

package locks_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/locks/tmpl"
)

const (
	testRegion        = "us-ord"
	testLocksDataName = "data.linode_locks.test"
)

func TestAccDataSourceLocks_basic(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLocksDataName, "locks.#", "1"),
					resource.TestCheckResourceAttrSet(testLocksDataName, "locks.0.id"),
					resource.TestCheckResourceAttrSet(testLocksDataName, "locks.0.entity_id"),
					resource.TestCheckResourceAttr(testLocksDataName, "locks.0.entity_type", "linode"),
					resource.TestCheckResourceAttr(testLocksDataName, "locks.0.lock_type", "cannot_delete"),
					resource.TestCheckResourceAttrSet(testLocksDataName, "locks.0.entity_label"),
					resource.TestCheckResourceAttrSet(testLocksDataName, "locks.0.entity_url"),
				),
			},
		},
	})
}

func TestAccDataSourceLocks_filter(t *testing.T) {
	t.Parallel()

	instanceLabel := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilter(t, instanceLabel, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testLocksDataName, "locks.#", "1"),
					resource.TestCheckResourceAttr(testLocksDataName, "locks.0.entity_type", "linode"),
					resource.TestCheckResourceAttr(testLocksDataName, "locks.0.lock_type", "cannot_delete"),
				),
			},
		},
	})
}
