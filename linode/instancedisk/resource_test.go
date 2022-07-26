package instancedisk_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancedisk/tmpl"
	"testing"
)

func TestAccResourceInstanceDisk_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_disk.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, 2048),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2048"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),

					resource.TestCheckResourceAttrSet(resName, "linode_id"),
				),
			},
			// Resize up
			{
				Config: tmpl.Basic(t, label, 2049),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2049"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),

					resource.TestCheckResourceAttrSet(resName, "linode_id"),
				),
			},
			// Resize down
			{
				Config: tmpl.Basic(t, label, 2047),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2047"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),

					resource.TestCheckResourceAttrSet(resName, "linode_id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceInstanceDisk_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_disk.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, label, 2048),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2048"),
					resource.TestCheckResourceAttr(resName, "filesystem", "ext4"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),
					resource.TestCheckResourceAttrSet(resName, "linode_id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkExists(name string, disk *linodego.InstanceDisk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		ids, err := helper.ParseMultiSegmentID(rs.Primary.ID, 2)
		if err != nil {
			return fmt.Errorf("failed to get disk info: %v", err)
		}

		linodeID, id := ids[0], ids[1]

		found, err := client.GetInstanceDisk(context.Background(), linodeID, id)
		if err != nil {
			return fmt.Errorf("error retrieving state of disk %s: %s", rs.Primary.Attributes["label"], err)
		}

		if disk != nil {
			*disk = *found
		}

		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance_disk" {
			continue
		}

		ids, err := helper.ParseMultiSegmentID(rs.Primary.ID, 2)
		if err != nil {
			return fmt.Errorf("failed to get disk info: %v", err)
		}

		linodeID, id := ids[0], ids[1]

		_, err = client.GetInstanceConfig(context.Background(), linodeID, id)

		if err == nil {
			return fmt.Errorf("disk with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting disk with id %d", id)
		}
	}

	return nil
}
