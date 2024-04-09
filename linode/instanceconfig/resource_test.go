//go:build integration || instanceconfig

package instanceconfig_test

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
	"github.com/linode/terraform-provider-linode/v2/linode/instanceconfig/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"vlans", "VPCs"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceInstanceConfig_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
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

func TestAccResourceInstanceConfig_deviceBlock(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	devicesCheck := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resName, "devices.0.sda.0.disk_id"),
		resource.TestCheckResourceAttrSet(resName, "devices.0.sdb.0.disk_id"),

		resource.TestCheckResourceAttrSet(resName, "device.0.disk_id"),
		resource.TestCheckResourceAttrSet(resName, "device.1.disk_id"),
		resource.TestCheckResourceAttr(resName, "device.0.device_name", "sda"),
		resource.TestCheckResourceAttr(resName, "device.1.device_name", "sdb"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			// Ensure the provider doesn't panic when creating an instance
			// with the new `device` block.
			{
				Config: tmpl.DeviceBlock(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					devicesCheck,
				),
			},
			{
				Config: tmpl.DeviceNamedBlock(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					devicesCheck,
				),
			},
			{
				Config: tmpl.DeviceBlock(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					devicesCheck,
				),
			},
			{
				Config: tmpl.DeviceNamedBlock(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					devicesCheck,
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

func TestAccResourceInstanceConfig_complex(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
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
				Config: tmpl.ComplexUpdates(t, instanceName, testRegion, rootPass, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config-updated"),
					resource.TestCheckResourceAttr(resName, "comments", "cool-updated"),

					resource.TestCheckResourceAttr(resName, "helpers.0.devtmpfs_automount", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.distro", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.modules_dep", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.network", "false"),
					resource.TestCheckResourceAttr(resName, "helpers.0.updatedb_disabled", "false"),

					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "interface.0.label", "cooler"),
					resource.TestCheckResourceAttr(resName, "interface.0.ipam_address", "10.0.0.3/24"),

					resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-32bit"),
					resource.TestCheckResourceAttr(resName, "memory_limit", "513"),
					resource.TestCheckResourceAttr(resName, "root_device", "/dev/sdb"),
					resource.TestCheckResourceAttr(resName, "virt_mode", "fullvirt"),

					resource.TestCheckResourceAttr(resName, "booted", "true"),

					resource.TestCheckResourceAttrSet(resName, "devices.0.sdb.0.disk_id"),
				),
			},
			{
				PreConfig: acceptance.AssertInstanceReboot(t, true, &instance),
				Config:    tmpl.ComplexUpdates(t, instanceName, testRegion, rootPass, true),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
				// Remove this ignorance when the TF SDK issue is fixed
				// https://github.com/hashicorp/terraform-plugin-sdk/issues/792
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccResourceInstanceConfig_booted(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Booted(t, instanceName, testRegion, false, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					resource.TestCheckResourceAttr(resName, "booted", "false"),
					resource.TestCheckResourceAttrSet(resName, "devices.0.sda.0.disk_id"),
				),
			},
			{
				PreConfig: func() {
					if instance.Status != linodego.InstanceOffline {
						t.Fatalf("expected instance to be offline, got %s", instance.Status)
					}
				},
				Config: tmpl.Booted(t, instanceName, testRegion, true, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					resource.TestCheckResourceAttr(resName, "booted", "true"),
					resource.TestCheckResourceAttrSet(resName, "devices.0.sda.0.disk_id"),
				),
			},
			{
				PreConfig: func() {
					if instance.Status != linodego.InstanceRunning {
						t.Fatalf("expected instance to be running, got %s", instance.Status)
					}
				},
				Config: tmpl.Booted(t, instanceName, testRegion, true, rootPass),
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

func TestAccResourceInstanceConfig_bootedSwap(t *testing.T) {
	t.Parallel()

	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		var instance linodego.Instance

		config1Name := "linode_instance_config.foobar1"
		config2Name := "linode_instance_config.foobar2"
		instanceName := acctest.RandomWithPrefix("tf_test")
		rootPass := acctest.RandString(12)

		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             checkDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.BootedSwap(t, instanceName, testRegion, false, rootPass),
					Check: resource.ComposeTestCheckFunc(
						acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
						checkExists(config1Name, nil),
						checkExists(config2Name, nil),

						resource.TestCheckResourceAttr(config1Name, "booted", "false"),
						resource.TestCheckResourceAttr(config2Name, "booted", "true"),
					),
				},
				{
					PreConfig: func() {
						if instance.Status != linodego.InstanceRunning {
							t.Fatalf("expected instance to be running, got %s", instance.Status)
						}
					},
					Config: tmpl.BootedSwap(t, instanceName, testRegion, true, rootPass),
					Check: resource.ComposeTestCheckFunc(
						acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
						checkExists(config1Name, nil),
						checkExists(config2Name, nil),

						resource.TestCheckResourceAttr(config1Name, "booted", "true"),
						resource.TestCheckResourceAttr(config2Name, "booted", "false"),
					),
				},
				{
					PreConfig: func() {
						if instance.Status != linodego.InstanceRunning {
							t.Fatalf("expected instance to be running, got %s", instance.Status)
						}
					},
					Config: tmpl.BootedSwap(t, instanceName, testRegion, true, rootPass),
				},
			},
		})
	})
}

func TestAccResourceInstanceConfig_provisioner(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	resName := "linode_instance_config.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Provisioner(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "label", "my-config"),
					resource.TestCheckResourceAttr(resName, "booted", "true"),
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

func TestAccResourceInstanceConfig_vpcInterface(t *testing.T) {
	t.Parallel()

	resName := "linode_instance_config.foobar"
	networkDSName := "data.linode_instance_networking.foobar"
	instanceName := acctest.RandomWithPrefix("tf-test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPCInterface(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "public"),
					resource.TestCheckResourceAttr(resName, "interface.1.purpose", "vpc"),
					resource.TestCheckResourceAttr(resName, "interface.1.ipv4.0.vpc", "10.0.4.250"),
					resource.TestCheckResourceAttr(resName, "interface.1.ip_ranges.0", "10.0.4.101/32"),
					resource.TestCheckResourceAttrSet(resName, "interface.1.ipv4.0.nat_1_1"),

					resource.TestCheckResourceAttr(networkDSName, "ipv4.0.public.0.vpc_nat_1_1.address", "10.0.4.250"),
					resource.TestCheckResourceAttrSet(networkDSName, "ipv4.0.public.0.vpc_nat_1_1.vpc_id"),
					resource.TestCheckResourceAttrSet(networkDSName, "ipv4.0.public.0.vpc_nat_1_1.subnet_id"),
				),
			},
			{
				Config: tmpl.VPCInterfaceUpdated(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "public"),
					resource.TestCheckResourceAttr(resName, "interface.1.purpose", "vpc"),
					resource.TestCheckResourceAttr(resName, "interface.1.ipv4.0.vpc", "10.0.4.249"),
					resource.TestCheckResourceAttr(resName, "interface.1.active", "false"),
					resource.TestCheckResourceAttr(resName, "interface.1.ip_ranges.0", "10.0.4.100/32"),

					resource.TestCheckResourceAttr(networkDSName, "ipv4.0.public.0.vpc_nat_1_1.#", "0"),
				),
			},
			{
				Config: tmpl.VPCInterfaceSwapped(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "interface.#", "2"),
					resource.TestCheckResourceAttr(resName, "interface.1.purpose", "public"),
					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "vpc"),
					resource.TestCheckResourceAttr(resName, "interface.0.ipv4.0.vpc", "10.0.4.249"),
					resource.TestCheckResourceAttr(resName, "interface.0.active", "false"),
					resource.TestCheckResourceAttr(resName, "interface.0.ip_ranges.0", "10.0.4.100/32"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: tmpl.VPCInterfaceOnly(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "vpc"),
					resource.TestCheckResourceAttr(resName, "interface.0.ipv4.0.vpc", "10.0.4.249"),
					resource.TestCheckResourceAttr(resName, "interface.0.active", "false"),
					resource.TestCheckResourceAttr(resName, "interface.0.ip_ranges.0", "10.0.4.100/32"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
			{
				Config: tmpl.VPCInterfaceRemoved(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					resource.TestCheckResourceAttr(resName, "interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "interface.0.purpose", "public"),
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

// Test case to ensure instances manually booted into rescue mode
// will not crash the provider.
func TestAccResourceInstanceConfig_rescueBooted(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	resName := "linode_instance_config.foobar"
	instanceResName := "linode_instance.foobar"

	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					acceptance.CheckInstanceExists(instanceResName, &instance),
				),
			},
			{
				PreConfig: func() {
					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatal(err)
					}

					poller, err := client.NewEventPoller(context.Background(), instance.ID, linodego.EntityLinode, linodego.ActionLinodeReboot)
					if err != nil {
						t.Fatalf("failed to create event poller: %v", err)
					}

					if err := client.RescueInstance(context.Background(), instance.ID, linodego.InstanceRescueOptions{}); err != nil {
						t.Fatalf("failed to boot instance into rescue mode: %v", err)
					}

					if _, err := poller.WaitForFinished(context.Background(), 240); err != nil {
						t.Fatalf("failed to wait for instance to boot into rescue mode: %v", err)
					}
				},
				Config: tmpl.ComplexUpdates(t, instanceName, testRegion, rootPass, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					acceptance.CheckInstanceExists(instanceResName, &instance),

					// The provider should not be booted into this config
					resource.TestCheckResourceAttr(resName, "booted", "false"),
				),
			},
			{
				Config: tmpl.Complex(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resName, nil),
					acceptance.CheckInstanceExists(instanceResName, &instance),

					// The provider should now have been rebooted into this config
					resource.TestCheckResourceAttr(resName, "booted", "true"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
				// Remove this ignorance when the TF SDK issue is fixed
				// https://github.com/hashicorp/terraform-plugin-sdk/issues/792
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func checkExists(name string, config *linodego.InstanceConfig) resource.TestCheckFunc {
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

		linodeID, id, err := getResourceIDs(rs)
		if err != nil {
			return fmt.Errorf("failed to get disk info: %v", err)
		}

		_, err = client.GetInstanceConfig(context.Background(), linodeID, id)

		if err == nil {
			return fmt.Errorf("config with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting config with id %d", id)
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
		if rs.Type != "linode_instance_config" {
			continue
		}

		linodeID, id, err := getResourceIDs(rs)
		if err != nil {
			return "", fmt.Errorf("failed to get config info: %v", err)
		}

		return fmt.Sprintf("%d,%d", linodeID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_instance_config")
}
