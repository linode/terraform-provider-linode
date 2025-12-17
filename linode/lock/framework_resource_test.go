//go:build integration || lock

package lock_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/lock/tmpl"
)

func TestAccResourceLock_basic(t *testing.T) {
	acceptance.OptInTest(t)
	t.Parallel()

	instanceName := "linode_instance.test"
	lockName := "linode_lock.test"

	var instance linodego.Instance

	label := acctest.RandomWithPrefix("tf_test")
	testRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories:  acceptance.ProtoV6ProviderFactories,
		CheckDestroy:              checkLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckInstanceExists(instanceName, &instance),
					resource.TestCheckResourceAttrSet(lockName, "id"),
					resource.TestCheckResourceAttrSet(lockName, "entity_id"),
					resource.TestCheckResourceAttr(lockName, "entity_type", "linode"),
					resource.TestCheckResourceAttr(lockName, "lock_type", "cannot_delete"),
					resource.TestCheckResourceAttrSet(lockName, "entity_label"),
					resource.TestCheckResourceAttrSet(lockName, "entity_url"),
				),
			},
			{
				ResourceName:      lockName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceLock_withSubresources(t *testing.T) {
	acceptance.OptInTest(t)
	t.Parallel()

	instanceName := "linode_instance.test"
	lockName := "linode_lock.test"

	var instance linodego.Instance

	label := acctest.RandomWithPrefix("tf_test")
	testRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories:  acceptance.ProtoV6ProviderFactories,
		CheckDestroy:              checkLockDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithSubresources(t, label, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckInstanceExists(instanceName, &instance),
					resource.TestCheckResourceAttrSet(lockName, "id"),
					resource.TestCheckResourceAttrSet(lockName, "entity_id"),
					resource.TestCheckResourceAttr(lockName, "entity_type", "linode"),
					resource.TestCheckResourceAttr(lockName, "lock_type", "cannot_delete_with_subresources"),
					resource.TestCheckResourceAttrSet(lockName, "entity_label"),
					resource.TestCheckResourceAttrSet(lockName, "entity_url"),
				),
			},
		},
	})
}

func checkLockDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lock" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetLock(context.Background(), id)
		if err == nil {
			return fmt.Errorf("Lock with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Lock with id %d: %s", id, err)
		}
	}

	return nil
}
