package instanceconfig_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instanceconfig/tmpl"
	"strconv"
	"strings"
	"testing"
)

func TestAccResourceInstanceConfig_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instanceName),
				Check: resource.ComposeTestCheckFunc(
					checkConfigExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
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

func TestAccResourceInstanceConfig_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, instanceName),
				Check: resource.ComposeTestCheckFunc(
					checkConfigExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					resource.TestCheckResourceAttr(resName, "comments", "cool"),

					resource.TestCheckResourceAttr(resName, "helpers.0.devtmpfs_automount", "true"),
					resource.TestCheckResourceAttr(resName, "helpers.0.distro", "true"),
					resource.TestCheckResourceAttr(resName, "helpers.0.modules_dep", "true"),
					resource.TestCheckResourceAttr(resName, "helpers.0.network", "true"),
					resource.TestCheckResourceAttr(resName, "helpers.0.updatedb_disabled", "true"),

					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "public"),

					resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "memory_limit", "512"),
					resource.TestCheckResourceAttr(resName, "root_device", "/dev/sda"),
					resource.TestCheckResourceAttr(resName, "virt_mode", "paravirt"),

					resource.TestCheckResourceAttr(resName, "booted", "true"),

					resource.TestCheckResourceAttrSet(resName, "devices.0.sda.0.disk_id"),
				),
			},
			{
				Config: tmpl.ComplexUpdates(t, instanceName),
				Check: resource.ComposeTestCheckFunc(
					checkConfigExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config-updated"),
					resource.TestCheckResourceAttr(resName, "comments", "cool-updated"),

					resource.TestCheckResourceAttr(resName, "helpers.0.devtmpfs_automount", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.distro", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.modules_dep", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.network", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.updatedb_disabled", "false"),

					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "interface.0.label", "cool"),
					resource.TestCheckResourceAttr(resName, "interface.0.ipam_address", "10.0.0.3/24"),

					resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-32bit"),
					resource.TestCheckResourceAttr(resName, "memory_limit", "513"),
					resource.TestCheckResourceAttr(resName, "root_device", "/dev/sdb"),
					resource.TestCheckResourceAttr(resName, "virt_mode", "fullvirt"),

					resource.TestCheckResourceAttr(resName, "booted", "false"),

					resource.TestCheckResourceAttrSet(resName, "devices.0.sdb.0.disk_id"),
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

func checkConfigExists(name string, config *linodego.InstanceConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		linodeID, id, err := getConfigInfo(rs)
		if err != nil {
			return fmt.Errorf("failed to get config info: %v", err)
		}

		found, err := client.GetInstanceConfig(context.Background(), linodeID, id)
		if err != nil {
			return fmt.Errorf("error retrieving state of config %s: %s", rs.Primary.Attributes["label"], err)
		}

		if config != nil {
			*config = *found
		}

		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_instance_config" {
			continue
		}

		linodeID, id, err := getConfigInfo(rs)
		if err != nil {
			return fmt.Errorf("failed to get config info: %v", err)
		}

		_, err = client.GetInstanceConfig(context.Background(), linodeID, id)

		if err == nil {
			return fmt.Errorf("config database with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting config with id %d", id)
		}
	}

	return nil
}

func getConfigInfo(rs *terraform.ResourceState) (int, int, error) {
	idSlug := strings.Split(rs.Primary.ID, "/")

	if len(idSlug) != 2 {
		return 0, 0, fmt.Errorf("invalid number of id segments")
	}

	linodeID, err := strconv.Atoi(idSlug[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
	}
	if linodeID == 0 {
		return 0, 0, fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, linodeID)
	}

	id, err := strconv.Atoi(idSlug[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
	}
	if id == 0 {
		return 0, 0, fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
	}

	return linodeID, id, nil
}
