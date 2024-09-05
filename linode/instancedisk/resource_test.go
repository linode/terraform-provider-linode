//go:build integration || instancedisk

package instancedisk_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instancedisk/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceInstanceDisk_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_disk.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, testRegion, 2048),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2048"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),

					resource.TestCheckResourceAttrSet(resName, "linode_id"),

					resource.TestCheckResourceAttrPair(
						resName, "disk_encryption",
						"linode_instance.foobar", "disk_encryption",
					),
				),
			},
			// Resize up
			{
				Config: tmpl.Basic(t, label, testRegion, 2049),
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
				Config: tmpl.Basic(t, label, testRegion, 2047),
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
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func TestAccResourceInstanceDisk_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_disk.foobar"
	label := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, label, testRegion, 2048),
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
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func TestAccResourceInstanceDisk_bootedResize(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_disk.foobar"
	label := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	var instance linodego.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.BootedResize(t, label, testRegion, 2048, rootPass),
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
				Config: tmpl.BootedResize(t, label, testRegion, 2049, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					resource.TestCheckResourceAttr(resName, "label", label),
					resource.TestCheckResourceAttr(resName, "size", "2049"),
					resource.TestCheckResourceAttr(resName, "status", "ready"),

					resource.TestCheckResourceAttrSet(resName, "linode_id"),
				),
			},
			{
				PreConfig: func() {
					if instance.Status != linodego.InstanceRunning {
						t.Fatalf("expected instance to be running, found %s", instance.Status)
					}
				},
				Config: tmpl.BootedResize(t, label, testRegion, 2049, rootPass),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       resourceImportStateID,
				ImportStateVerifyIgnore: []string{"image", "root_pass"},
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

		linodeID, id, err := getResourceIDs(rs)
		if err != nil {
			return fmt.Errorf("failed to get disk info: %v", err)
		}

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

		linodeID, id, err := getResourceIDs(rs)
		if err != nil {
			return fmt.Errorf("failed to get disk info: %v", err)
		}

		_, err = client.GetInstanceDisk(context.Background(), linodeID, id)

		if err == nil {
			return fmt.Errorf("disk with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting disk with id %d", id)
		}
	}

	return nil
}

func getResourceIDs(rs *terraform.ResourceState) (int, int, error) {
	id, err := strconv.Atoi(rs.Primary.ID)
	if err != nil {
		return 0, 0, err
	}

	linodeID, err := strconv.Atoi(rs.Primary.Attributes["linode_id"])
	if err != nil {
		return 0, 0, err
	}

	return linodeID, id, nil
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance_disk" {
			continue
		}

		linodeID, id, err := getResourceIDs(rs)
		if err != nil {
			return "", fmt.Errorf("failed to get disk info: %v", err)
		}

		return fmt.Sprintf("%d,%d", linodeID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_instance_disk")
}
