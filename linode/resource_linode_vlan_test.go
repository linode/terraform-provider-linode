package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/linode/linodego"
)

const (
	testVLANResName     = "linode_vlan.test"
	testVLANDescription = "terraform-provider-linode acctest"
)

func init() {
	resource.AddTestSweepers("linode_vlan", &resource.Sweeper{
		Name: "linode_vlan",
		F:    testSweepLinodeVLAN,
	})
}

func testSweepLinodeVLAN(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err)
	}

	vlan, err := client.ListVLANs(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get VLANs: %s", err)
	}
	for _, vlan := range vlan {
		if vlan.Description != testVLANDescription {
			continue
		}
		if err := client.DeleteVLAN(context.Background(), vlan.ID); err != nil {
			return fmt.Errorf("failed to destroy VLAN %d during sweep: %s", vlan.ID, err)
		}
	}

	return nil
}

func testAccCheckLinodeVLANAttachedLinodes(linodes ...*linodego.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(linodego.Client)

		rs, ok := s.RootModule().Resources[testVLANResName]
		if !ok {
			return fmt.Errorf("could not find: %s", testVLANResName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed parsing %v to int", rs.Primary.ID)
		}

		vlan, err := client.GetVLAN(context.Background(), id)
		if err != nil {
			return fmt.Errorf("failed to find VLAN %d: %s", id, err)
		}

		attachedLinodes := make(map[int]struct{}, len(vlan.Linodes))
		for _, linode := range vlan.Linodes {
			attachedLinodes[linode.ID] = struct{}{}
		}

		for _, linode := range linodes {
			if _, ok := attachedLinodes[linode.ID]; !ok {
				return fmt.Errorf("expected linode %d to be attached to vlan %d", linode.ID, vlan.ID)
			} else {
				delete(attachedLinodes, linode.ID)
			}
		}

		if len(attachedLinodes) != 0 {
			return fmt.Errorf("unexpected linodes %v attached to vlan %d", linodes, vlan.ID)
		}
		return nil
	}
}

func testAccCheckLinodeVLANDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	rs, ok := s.RootModule().Resources[testVLANResName]
	if !ok {
		return fmt.Errorf("could not find: %s", testVLANResName)
	}

	id, err := strconv.Atoi(rs.Primary.ID)
	if err != nil {
		return fmt.Errorf("failed to parse LKE Cluster ID: %s", err)
	}

	if id == 0 {
		return fmt.Errorf("should not have LKE Cluster ID of 0")
	}

	if _, err = client.GetVLAN(context.Background(), id); err == nil {
		return fmt.Errorf("expected VLAN %d to have been deleted", id)
	} else if apiErr, ok := err.(*linodego.Error); !ok {
		return fmt.Errorf("expected API Error but got %#v", err)
	} else if apiErr.Code != 404 {
		return fmt.Errorf("expected an error 404 but got %#v", apiErr)
	}
	return nil
}

func TestAccLinodeVLAN_basic(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVLANBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.test", &instance),
					resource.TestCheckResourceAttr(testVLANResName, "region", "ca-central"),
					resource.TestCheckResourceAttr(testVLANResName, "description", testVLANDescription),
					resource.TestCheckResourceAttr(testVLANResName, "cidr_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testVLANResName, "attached_linodes.#", "1"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.mac_address"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.ipv4_address"),
					testAccCheckLinodeVLANAttachedLinodes(&instance),
				),
			},
		},
	})
}

func TestAccLinodeVLAN_updateMultipleLinodes(t *testing.T) {
	t.Parallel()

	var instance1, instance2, instance3 linodego.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVLANBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.test", &instance1),
					resource.TestCheckResourceAttr(testVLANResName, "attached_linodes.#", "1"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.mac_address"),
					testAccCheckLinodeVLANAttachedLinodes(&instance1),
				),
			},
			{
				Config: testAccCheckLinodeVLANWithMultipleLinodes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeInstanceExists("linode_instance.test2", &instance2),
					testAccCheckLinodeInstanceExists("linode_instance.test3", &instance3),
					resource.TestCheckResourceAttr(testVLANResName, "cidr_block", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(testVLANResName, "attached_linodes.#", "2"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.0.mac_address"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.1.ipv4_address"),
					resource.TestCheckResourceAttrSet(testVLANResName, "attached_linodes.1.mac_address"),
					testAccCheckLinodeVLANAttachedLinodes(&instance2, &instance3),
				),
			},
		},
	})
}

func testAccCheckLinodeVLANBasic() string {
	return testAccCheckLinodeInstanceWithBootImage("test", acctest.RandomWithPrefix("tf_test")) +
		fmt.Sprintf(`
resource "linode_vlan" "test" {
	description = "%s"
	region      = "ca-central"
	linodes     = [linode_instance.test.id]
	cidr_block  = "0.0.0.0/0"
}`, testVLANDescription)
}

func testAccCheckLinodeVLANWithMultipleLinodes() string {
	return testAccCheckLinodeInstanceWithBootImage("test", acctest.RandomWithPrefix("tf_test")) +
		testAccCheckLinodeInstanceWithBootImage("test2", acctest.RandomWithPrefix("tf_test")) +
		testAccCheckLinodeInstanceWithBootImage("test3", acctest.RandomWithPrefix("tf_test")) +
		fmt.Sprintf(`
resource "linode_vlan" "test" {
	description = "%s"
	region      = "ca-central"
	linodes     = [linode_instance.test2.id, linode_instance.test3.id]
	cidr_block  = "0.0.0.0/0"
}`, testVLANDescription)
}
