//go:build integration || instance

package instance_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instance/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_instance", &resource.Sweeper{
		Name: "linode_instance",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Vlans", "VPCs"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	instances, err := client.ListInstances(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting instances: %s", err)
	}
	for _, instance := range instances {
		if !acceptance.ShouldSweep(prefix, instance.Label) {
			continue
		}
		err := client.DeleteInstance(context.Background(), instance.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", instance.Label, err)
		}
	}

	return nil
}

func TestAccResourceInstance_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "256"),
					resource.TestCheckResourceAttrSet(resName, "host_uuid"),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_watchdogDisabled(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WatchdogDisabled(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "watchdog_enabled", "false"),
				),
			},
			{
				Config:   tmpl.WatchdogDisabled(t, instanceName, testRegion, rootPass),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceInstance_authorizedUsers(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.AuthorizedUsers(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "256"),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_users", "image", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_validateAuthorizedKeys(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.AuthorizedKeysEmpty(t, instanceName, testRegion),
				ExpectError: regexp.MustCompile(
					"invalid input for authorized_keys"),
			},
			{
				Config: tmpl.DiskAuthorizedKeysEmpty(t, instanceName, testRegion, rootPass),
				ExpectError: regexp.MustCompile(
					"invalid input for disk authorized_keys"),
			},
		},
	})
}

func TestAccResourceInstance_interfaces(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Interfaces(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),

					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),

					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.label", "tf-really-cool-vlan"),
				),
			},
			{
				Config: tmpl.InterfacesUpdate(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "2"),

					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "public"),

					resource.TestCheckResourceAttr(resName, "config.0.interface.1.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.1.label", "tf-really-cool-vlan"),
				),
			},
			{
				Config: tmpl.InterfacesUpdateEmpty(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image", "interface", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_config(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithConfig(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "alerts.0.cpu", "60"),

					resource.TestCheckResourceAttrSet(resName, "config.0.id"),
					resource.TestCheckResourceAttr(resName, "config.0.run_level", "binbash"),
					resource.TestCheckResourceAttr(resName, "config.0.virt_mode", "fullvirt"),
					resource.TestCheckResourceAttr(resName, "config.0.memory_limit", "1024"),

					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_configPair(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.MultipleConfigs(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkComputeInstanceConfigs(&instance, testConfig("configa", testConfigKernel("linode/latest-64bit"))),
					checkComputeInstanceConfigs(&instance, testConfig("configb", testConfigKernel("linode/latest-32bit"))),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"boot_config_label", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_configInterfaces(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ConfigInterfaces(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),

					resource.TestCheckResourceAttr(resName, "config.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.label", "tf-really-cool-vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
				),
			},
			{
				Config: tmpl.ConfigInterfacesMultiple(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.label", "tf-really-cool-vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "config.1.interface.#", "2"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
				),
			},
			{
				PreConfig: testAccAssertReboot(t, true, &instance),
				Config:    tmpl.ConfigInterfacesUpdate(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
				),
			},
			{
				PreConfig: testAccAssertReboot(t, true, &instance),
				Config:    tmpl.ConfigInterfacesUpdateEmpty(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "0"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_configInterfacesNoReboot(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ConfigInterfaces(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),

					resource.TestCheckResourceAttr(resName, "config.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.label", "tf-really-cool-vlan"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
				),
			},
			{
				Config: tmpl.ConfigInterfacesUpdateNoReboot(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "config.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "0"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
				),
			},
			{
				PreConfig: testAccAssertReboot(t, false, &instance),
				Config:    tmpl.ConfigInterfacesUpdateNoReboot(t, instanceName, testRegion, rootPass),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "boot_config_label", "migration_type"},
			},
		},
	})
}

func testAccAssertReboot(t *testing.T, shouldRestart bool, instance *linodego.Instance) func() {
	return func() {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
		eventFilter := fmt.Sprintf(`{"entity.type": "linode", "entity.id": %d, "action": "linode_reboot", "created": { "+gte": "%s" }}`,
			instance.ID, instance.Created.Format("2006-01-02T15:04:05"))

		events, err := client.ListEvents(context.Background(), &linodego.ListOptions{Filter: eventFilter})
		if err != nil {
			t.Fail()
		}

		if len(events) == 0 && shouldRestart {
			t.Fatal("expected instance to have been rebooted")
		}

		if len(events) > 0 && !shouldRestart {
			t.Fatal("expected instance to not have been rebooted")
		}
	}
}

func TestAccResourceInstance_disk(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.RawDisk(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
					resource.TestCheckResourceAttr(resName, "config.#", "0"),
					resource.TestCheckResourceAttr(resName, "disk.#", "1"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					resource.TestCheckResourceAttr(resName, "disk.0.label", "disk"),
					checkComputeInstanceDisk(&instance, "disk", 3000),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_diskImage(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.Disk(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					checkComputeInstanceDisk(&instance, "disk", 3000),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_diskPair(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	var instanceDisk linodego.InstanceDisk
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DiskMultiple(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "512"),
					checkInstanceDisks(&instance,
						testDisk("diska", testDiskSize(3000), testDiskExists(&instanceDisk)),
						testDisk("diskb", testDiskSize(512)),
					),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_diskAndConfig(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkComputeInstanceConfigs(&instance,
						testConfig("config", testConfigKernel("linode/latest-64bit")),
					),
					checkComputeInstanceDisk(&instance, "disk", 3000),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_disksAndConfigs(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	var instanceDisk linodego.InstanceDisk

	instanceName := acctest.RandomWithPrefix("tf_test")

	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			acceptance.CheckInstanceDestroy,
			acceptance.CheckVolumeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: tmpl.DiskConfigMultiple(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "512"),
					checkInstanceDiskExists(&instance, "diska", &instanceDisk),
					// TODO(displague) create checkInstanceDisks helper (like Configs)
					checkComputeInstanceDisk(&instance, "diska", 3000),
					checkComputeInstanceDisk(&instance, "diskb", 512),
					checkComputeInstanceConfigs(&instance,
						testConfig("configa", testConfigKernel("linode/latest-64bit"), testConfigSDADisk(&instanceDisk)),
						testConfig("configb", testConfigKernel("linode/grub2"), testConfigComments("won't boot"), testConfigSDBDisk(&instanceDisk)),
					),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"boot_config_label", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_volumeAndConfig(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	volName := "linode_volume.foo"

	var instance linodego.Instance
	var instanceDisk linodego.InstanceDisk
	var volume linodego.Volume
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VolumeConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					acceptance.CheckVolumeExists(volName, &volume),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "boot_config_label", "config"),
					checkInstanceDiskExists(&instance, "disk", &instanceDisk),
					// TODO(displague) create checkInstanceDisks helper (like Configs)
					checkComputeInstanceDisk(&instance, "disk", 3000),
					checkComputeInstanceConfigs(&instance,
						testConfig("config", testConfigKernel("linode/latest-64bit"), testConfigSDADisk(&instanceDisk), testConfigSDBVolume(&volume)),
					),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_privateImage(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.PrivateImage(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					checkInstanceDisks(&instance,
						testDisk("boot", testDiskSize(1000)),
						testDisk("swap", testDiskSize(800)),
						testDisk("logs", testDiskSize(600)),
					),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_noImage(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.NoImage(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_updateSimple(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
				),
			},
			{
				Config: tmpl.Updates(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", instanceName)),
					resource.TestCheckResourceAttr(resName, "group", "tf_test_r"),
				),
			},
		},
	})
}

func TestAccResourceInstance_configUpdate(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"

	// This test can occasionally fail while running the entire test suite in parallel
	acceptance.RunTestRetry(t, 3, func(retryT *acceptance.TRetry) {
		resource.Test(retryT, resource.TestCase{
			PreCheck:                 func() { acceptance.PreCheck(t) },
			ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
			CheckDestroy:             acceptance.CheckInstanceDestroy,
			Steps: []resource.TestStep{
				{
					Config: tmpl.WithConfig(t, instanceName, testRegion),
					Check: resource.ComposeTestCheckFunc(
						acceptance.CheckInstanceExists(resName, &instance),
						resource.TestCheckResourceAttr(resName, "label", instanceName),
						resource.TestCheckResourceAttr(resName, "group", "tf_test"),
						resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-64bit"),
						resource.TestCheckResourceAttr(resName, "config.0.root_device", "/dev/sda"),
						resource.TestCheckResourceAttr(resName, "config.0.helpers.0.network", "true"),
						resource.TestCheckResourceAttr(resName, "alerts.0.cpu", "60"),
					),
				},
				{
					Config: tmpl.ConfigUpdates(t, instanceName, testRegion),
					Check: resource.ComposeTestCheckFunc(
						acceptance.CheckInstanceExists(resName, &instance),
						resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", instanceName)),
						resource.TestCheckResourceAttr(resName, "group", "tf_test_r"),
						// changed kernel, not label
						resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
						resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-32bit"),
						resource.TestCheckResourceAttr(resName, "config.0.root_device", "/dev/sda"),
						resource.TestCheckResourceAttr(resName, "config.0.helpers.0.network", "false"),
						resource.TestCheckResourceAttr(resName, "alerts.0.cpu", "80"),
					),
				},
			},
		})
	})
}

func TestAccResourceInstance_configPairUpdate(t *testing.T) {
	t.Parallel()

	config := linodego.InstanceConfig{}
	configA := linodego.InstanceConfig{}
	configB := linodego.InstanceConfig{}

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithConfig(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "config.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-64bit"),
					checkComputeInstanceConfigs(&instance,
						testConfig("config", testConfigExists(&config), testConfigKernel("linode/latest-64bit")),
					),
				),
			},
			{
				Config: tmpl.MultipleConfigs(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "config.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "configa"),
					resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "config.1.label", "configb"),
					resource.TestCheckResourceAttr(resName, "config.1.kernel", "linode/latest-32bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkComputeInstanceConfigs(&instance,
						testConfig("configa", testConfigExists(&configA), testConfigKernel("linode/latest-64bit")),
						testConfig("configb", testConfigExists(&configB), testConfigKernel("linode/latest-32bit")),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"boot_config_label", "status", "resize_disk", "migration_type"},
			},
			{
				Config: tmpl.WithConfig(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "config.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-64bit"),
					checkComputeInstanceConfigs(&instance,
						testConfig("config", testConfigExists(&config), testConfigKernel("linode/latest-64bit")),
					),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"boot_config_label", "status", "resize_disk", "migration_type"},
			},
			{
				Config: tmpl.ConfigsAllUpdated(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					// resource.TestCheckResourceAttr(resName, "kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkComputeInstanceConfigs(&instance,
						testConfig("configb", testConfigKernel("linode/latest-64bit")),
						testConfig("configa", testConfigKernel("linode/latest-32bit")),
						testConfig("configc", testConfigKernel("linode/latest-64bit")),
					),
				),
			},
		},
	})
}

func TestAccResourceInstance_upsizeWithoutDisk(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithType(t, instanceName, acceptance.PublicKeyMaterial, "g6-nanode-1", testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25344)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
			{
				Config: tmpl.WithType(t, instanceName, acceptance.PublicKeyMaterial, "g6-standard-1", testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25344)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
		},
	})
}

func TestAccResourceInstance_diskRawResize(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			{
				Config: tmpl.RawDisk(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "config.#", "0"),
					resource.TestCheckResourceAttr(resName, "disk.#", "1"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					resource.TestCheckResourceAttr(resName, "disk.0.label", "disk"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(3000))),
				),
			},
			// Bump it to a 2048, and expand the disk
			{
				Config: tmpl.RawDiskExpanded(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					resource.TestCheckResourceAttr(resName, "config.#", "0"),
					resource.TestCheckResourceAttr(resName, "disk.#", "1"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "6000"),
					resource.TestCheckResourceAttr(resName, "disk.0.label", "disk"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(6000))),
				),
			},
		},
	})
}

func TestAccResourceInstance_tag(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a single tag
			{
				Config: tmpl.Tag(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
				),
			},
			// Apply updated tags
			{
				Config: tmpl.TagUpdate(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
			// Reapply with different case, expect no planned changes
			{
				Config:   tmpl.TagUpdateCaseChange(t, instanceName, testRegion),
				PlanOnly: true,
			},
			// Update the tags again, expect changes
			{
				Config: tmpl.Tag(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
				),
			},
		},
	})
}

func TestAccResourceInstance_tagWithVolume(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance

	label := acctest.RandomWithPrefix("tf_test")

	instanceResName := "linode_instance.foobar"
	volumeResName := "linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.TagVolume(t, label, "tf_test", testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(instanceResName, &instance),
					resource.TestCheckResourceAttr(instanceResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceResName, "tags.0", "tf_test"),
				),
			},
			{
				Config: tmpl.TagVolume(t, label, "tf_test_updated", testRegion),
				Check: resource.ComposeTestCheckFunc(
					// Ensure the volume is not detached
					acceptance.CheckEventAbsent(volumeResName, "volume", linodego.ActionVolumeDetach),

					acceptance.CheckInstanceExists(instanceResName, &instance),
					resource.TestCheckResourceAttr(instanceResName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceResName, "tags.0", "tf_test_updated"),
				),
			},
		},
	})
}

func TestAccResourceInstance_diskResize(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(3000))),
				),
			},
			// Increase disk size
			{
				Config: tmpl.DiskConfigResized(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "6000"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(6000))),
				),
			},
		},
	})
}

func TestAccResourceInstance_withDiskLinodeUpsize(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start with g6-nanode-1
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(3000))),
				),
			},
			// Upsize to g6-standard-1 with fully allocated disk
			{
				Config: tmpl.DiskConfigExpanded(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "51200"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(51200))),
				),
			},
		},
	})
}

func TestAccResourceInstance_withDiskLinodeDownsize(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start with g6-standard-1 with fully allocated disk
			{
				Config: tmpl.DiskConfigExpanded(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "51200"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(51200))),
				),
			},
			// Downsize to g6-nanode-1
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(3000))),
				),
			},
		},
	})
}

func TestAccResourceInstance_downsizeWithoutDisk(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.WithType(t, instanceName, acceptance.PublicKeyMaterial, "g6-standard-1", testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(50944)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
			{
				Config: tmpl.WithType(t, instanceName, acceptance.PublicKeyMaterial, "g6-nanode-1", testRegion, rootPass),
				ExpectError: regexp.MustCompile(
					"insufficient disk capacity"),
			},
		},
	})
}

func TestAccResourceInstance_fullDiskSwapUpsize(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	stackScriptName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.FullDisk(t, instanceName, acceptance.PublicKeyMaterial, stackScriptName, testRegion, 256, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25344)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
			{
				PreConfig: func() {
					ctx := context.Background()
					client := acceptance.GetSSHClient(t, "root", instance.IPv4[0].String())

					defer client.Close()
					ss, err := client.NewSession()
					if err != nil {
						t.Fatalf("failed to establish SSH session: %s", err)
					}

					ctx, cancel := context.WithTimeout(ctx, time.Minute)
					defer cancel()

					ticker := time.NewTicker(500 * time.Millisecond)
					defer ticker.Stop()

					for {
						select {
						case <-ticker.C:
							buf := new(bytes.Buffer)
							ss.Stdout = buf
							ss.Run("[[ $(df /dev/sda --block-size=1 | tail -n-1 | awk '{print $5}') == '100%' ]] && echo 1 || echo 0")

							if buf.String() == "1" {
								return
							}

						case <-ctx.Done():
							return
						}
					}
				},
				Config:      tmpl.FullDisk(t, instanceName, acceptance.PublicKeyMaterial, stackScriptName, testRegion, 512, rootPass),
				ExpectError: regexp.MustCompile("Error waiting for resize of Instance \\d+ Disk \\d+"),
			},
		},
	})
}

func TestAccResourceInstance_swapUpsize(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithSwapSize(t, instanceName, acceptance.PublicKeyMaterial, testRegion, 256, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25344)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
			{
				Config: tmpl.WithSwapSize(t, instanceName, acceptance.PublicKeyMaterial, testRegion, 512, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25088)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(512)),
					),
				),
			},
		},
	})
}

func TestAccResourceInstance_swapDownsize(t *testing.T) {
	t.Parallel()

	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"

	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.WithSwapSize(t, instanceName, acceptance.PublicKeyMaterial, testRegion, 512, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25088)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(512)),
					),
				),
			},
			{
				Config: tmpl.WithSwapSize(t, instanceName, acceptance.PublicKeyMaterial, testRegion, 256, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstanceDisks(&instance,
						testDiskByFS(linodego.FilesystemExt4, testDiskSize(25344)),
						testDiskByFS(linodego.FilesystemSwap, testDiskSize(256)),
					),
				),
			},
		},
	})
}

func TestAccResourceInstance_diskResizeAndExpanded(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(3000))),
				),
			},

			// Bump to 2048 and expand disk
			{
				Config: tmpl.DiskConfigResizedExpanded(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),

					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "6000"),

					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
					checkInstanceDisks(&instance, testDisk("disk", testDiskSize(6000))),
				),
			},
		},
	})
}

func TestAccResourceInstance_diskSlotReorder(t *testing.T) {
	t.Parallel()
	var (
		instance     linodego.Instance
		instanceDisk linodego.InstanceDisk
	)
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Start off with a Linode 1024
			{
				Config: tmpl.DiskConfig(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "25600"),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					checkInstanceDisks(&instance, testDisk("disk", testDiskExists(&instanceDisk), testDiskSize(3000))),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"), testConfigSDADisk(&instanceDisk))),
					resource.TestCheckResourceAttrSet(resName, "config.0.devices.0.sda.0.disk_id"),
					resource.TestCheckResourceAttr(resName, "config.0.devices.0.sdb.#", "0"),
					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					checkComputeInstanceConfigs(&instance, testConfig("config", testConfigKernel("linode/latest-64bit"))),
				),
			},
			// Add a disk, reorder the disks
			{
				Config: tmpl.DiskConfigReordered(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeAggregateTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "specs.0.disk", "51200"),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
					resource.TestCheckResourceAttr(resName, "disk.0.size", "3000"),
					resource.TestCheckResourceAttr(resName, "disk.0.label", "disk"),
					resource.TestCheckResourceAttrSet(resName, "disk.0.id"),
					resource.TestCheckResourceAttr(resName, "disk.1.size", "3000"),
					resource.TestCheckResourceAttr(resName, "disk.1.label", "diskb"),
					resource.TestCheckResourceAttrSet(resName, "disk.1.id"),
					resource.TestCheckResourceAttr(resName, "config.0.label", "config"),
					resource.TestCheckResourceAttr(resName, "config.0.kernel", "linode/latest-64bit"),
					resource.TestCheckResourceAttrSet(resName, "config.0.devices.0.sda.0.disk_id"),
					resource.TestCheckResourceAttrSet(resName, "config.0.devices.0.sdb.0.disk_id"),
					resource.TestCheckResourceAttr(resName, "config.0.devices.0.sdc.#", "0"),
					resource.TestCheckResourceAttrPair(resName, "config.0.devices.0.sda.0.disk_id", resName, "disk.1.id"),
					resource.TestCheckResourceAttrPair(resName, "config.0.devices.0.sdb.0.disk_id", resName, "disk.0.id"),

					resource.TestCheckResourceAttr(resName, "swap_size", "0"),
					resource.TestCheckResourceAttr(resName, "status", "running"),
				),
			},
		},
	})
}

func TestAccResourceInstance_privateNetworking(t *testing.T) {
	t.Parallel()
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_instance.foobar"
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.PrivateNetworking(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					checkInstancePrivateNetworkAttributes("linode_instance.foobar"),
					resource.TestCheckResourceAttr(resName, "private_ip", "true"),
				),
			},
		},
	})
}

func TestAccResourceInstance_stackScriptInstance(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.StackScript(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "group", "tf_test"),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_diskImageUpdate(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DiskBootImage(t, instanceName, acceptance.TestImagePrevious, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName)),
			},
			{
				Config: tmpl.DiskBootImage(t, instanceName, acceptance.TestImageLatest, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					// resource was tainted for recreation due to change of disk.0.image, marked
					// with ForceNew.
					acceptance.CheckResourceAttrNotEqual(resName, "id", strconv.Itoa(instance.ID)),
				),
			},

			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_stackScriptDisk(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DiskStackScript(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
				),
			},
		},
	})
}

func TestAccResourceInstance_typeChangeDiskImplicit(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"

	var instance linodego.Instance
	// oldDiskSize := 0

	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Create an initial instance
			{
				Config: tmpl.TypeChangeDisk(t, instanceName, "g6-nanode-1", testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
			// Upsize the instance and disk
			{
				Config: tmpl.TypeChangeDisk(t, instanceName, "g6-standard-1", testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
				),
			},
			// Attempt a downsize
			{
				Config:      tmpl.TypeChangeDisk(t, instanceName, "g6-nanode-1", testRegion, true),
				ExpectError: regexp.MustCompile("Did you try to resize a linode with implicit"),
			},
		},
	})
}

func TestAccResourceInstance_typeChangeDiskExplicit(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Create an instance with explicit disks
			{
				Config: tmpl.TypeChangeDiskExplicit(t, instanceName, "g6-nanode-1", testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
			// Attempt to resize the instance and disk and expect an error
			{
				Config:      tmpl.TypeChangeDiskExplicit(t, instanceName, "g6-standard-1", testRegion, true),
				ExpectError: regexp.MustCompile("all of `image,resize_disk` must be specified"),
			},
			// Resize only the instance
			{
				Config: tmpl.TypeChangeDiskExplicit(t, instanceName, "g6-standard-1", testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
				),
			},
		},
	})
}

func TestAccResourceInstance_typeChangeNoDisks(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			// Create an instance with explicit disks
			{
				Config: tmpl.TypeChangeDiskNone(t, instanceName, "g6-nanode-1", testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
			// Attempt to resize the instance
			{
				Config: tmpl.TypeChangeDiskNone(t, instanceName, "g6-standard-1", testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-standard-1"),
				),
			},
			// Attempt to downsize the instance
			{
				Config: tmpl.TypeChangeDiskNone(t, instanceName, "g6-nanode-1", testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
				),
			},
		},
	})
}

func TestAccResourceInstance_powerStateUpdates(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.BootState(t, instanceName, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
				),
			},
			{
				Config: tmpl.BootState(t, instanceName, testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "running"),
				),
			},
			{
				Config: tmpl.BootState(t, instanceName, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
				),
			},
			// Ensure an implicit reboot isn't triggered when booted == false
			{
				Config: tmpl.BootStateInterface(t, instanceName, testRegion, false),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
				),
			},
			{
				PreConfig: func() {
					client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

					filter := &linodego.Filter{}
					filter.AddField(linodego.Eq, "action", linodego.ActionLinodeReboot)
					filter.AddField(linodego.Eq, "entity.id", instance.ID)
					filter.AddField(linodego.Eq, "entity.type", linodego.EntityLinode)
					jsonData, err := filter.MarshalJSON()
					if err != nil {
						t.Fatal(err)
					}

					events, err := client.ListEvents(context.Background(), &linodego.ListOptions{Filter: string(jsonData)})
					if err != nil {
						t.Fatal(err)
					}

					if len(events) > 0 {
						t.Fatal("found reboot event when no reboot was expected")
					}
				},
				Config: tmpl.BootStateInterface(t, instanceName, testRegion, false),
			},
		},
	})
}

func TestAccResourceInstance_powerStateConfigUpdates(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.BootStateConfig(t, instanceName, testRegion, false, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
				),
			},
			{
				Config: tmpl.BootStateConfig(t, instanceName, testRegion, true, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "running"),
				),
			},
			{
				Config: tmpl.BootStateConfig(t, instanceName, testRegion, false, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "offline"),
				),
			},
		},
	})
}

func TestAccResourceInstance_powerStateConfigBooted(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.BootStateConfig(t, instanceName, testRegion, true, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "running"),
				),
			},
		},
	})
}

func TestAccResourceInstance_powerStateBooted(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.BootState(t, instanceName, testRegion, true),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "status", "running"),
				),
			},
		},
	})
}

func TestAccResourceInstance_powerStateNoImage(t *testing.T) {
	t.Parallel()

	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.BootStateNoImage(t, instanceName, testRegion, true),
				ExpectError: regexp.MustCompile("booted requires an image or disk/config be defined"),
			},
		},
	})
}

func TestAccResourceInstance_ipv4Sharing(t *testing.T) {
	t.Parallel()

	// We need to manually override the region as IP sharing capabilities aren't
	// explicitly mentioned by the API.
	const region = "us-west"

	failoverResName := "linode_instance.failover"

	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      tmpl.IPv4SharingBadInput(t, instanceName, region),
				ExpectError: regexp.MustCompile("expected ipv4 address, got"),
			},
			{
				Config: tmpl.IPv4Sharing(t, instanceName, region),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(failoverResName, &instance),
					resource.TestCheckResourceAttr(failoverResName, "shared_ipv4.#", "1"),
					resource.TestCheckResourceAttrSet(failoverResName, "shared_ipv4.0"),
				),
			},
			{
				Config: tmpl.IPv4SharingAllocation(t, instanceName, region),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(failoverResName, &instance),
					resource.TestCheckResourceAttr(failoverResName, "shared_ipv4.#", "1"),
					resource.TestCheckResourceAttrSet(failoverResName, "shared_ipv4.0"),
				),
			},
			{
				Config: tmpl.IPv4SharingEmpty(t, instanceName, region),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(failoverResName, &instance),
					resource.TestCheckResourceAttr(failoverResName, "shared_ipv4.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceInstance_userData(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Metadata"})
	if err != nil {
		t.Fatal(err)
	}

	rootPass := acctest.RandString(12)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.UserData(t, instanceName, region, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", region),

					resource.TestCheckResourceAttr(resName, "has_user_data", "true"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "metadata", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_requestQuantity(t *testing.T) {
	t.Parallel()

	const maxRequestsPerSecond = 3.0

	// We need to make sure we're not running into a race condition here
	var numRequestsLock sync.Mutex
	numRequests := 0
	var startTime time.Time

	instanceName := acctest.RandomWithPrefix("tf_test")

	provider, providerMap := acceptance.CreateTestProvider()

	rootPass := acctest.RandString(12)

	acceptance.ModifyProviderMeta(provider,
		func(ctx context.Context, config *helper.ProviderMeta) error {
			config.Client.OnBeforeRequest(func(request *linodego.Request) error {
				if startTime.IsZero() {
					startTime = time.Now()
				}

				numRequestsLock.Lock()
				defer numRequestsLock.Unlock()
				numRequests++

				return nil
			})

			return nil
		})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: providerMap,
		Steps: []resource.TestStep{
			{
				// Provision a bunch of Linodes and wait for them to boot into an image
				Config: tmpl.ManyLinodes(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
			},
			{
				PreConfig: func() {
					requestsPerSecond := (float64(numRequests) / float64(time.Since(startTime).Seconds()))

					t.Logf("\n[INFO] results from 12 linode parallel creation:\n"+
						"total requests: %d\nfrequency: ~%f requests/second\n", numRequests, requestsPerSecond)

					if requestsPerSecond > maxRequestsPerSecond {
						t.Fatalf("too many requests: %f > %f", requestsPerSecond, maxRequestsPerSecond)
					}
				},
				Config: tmpl.ManyLinodes(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
			},
		},
	})
}

func TestAccResourceInstance_firewallOnCreation(t *testing.T) {
	t.Parallel()

	instanceResourceName := "linode_instance.foobar"
	firewallResourceName := "linode_firewall.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Cloud Firewall"})
	rootPass := acctest.RandString(12)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.FirewallOnCreation(t, instanceName, region, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(instanceResourceName, &instance),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						firewallResourceName,
						"devices.0.label",
						instanceResourceName,
						"label",
					),
				),
			},
			{
				ResourceName:            instanceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "firewall_id", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_VPCInterface(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.VPCInterface(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vpc"),

					resource.TestCheckResourceAttr(resName, "config.0.interface.0.ipv4.0.vpc", "10.0.4.150"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.ip_ranges.0", "10.0.4.100/32"),
					resource.TestCheckResourceAttrSet(resName, "config.0.interface.0.ipv4.0.nat_1_1"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image", "interface", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_VPCPublicInterfacesAddRemoveSwap(t *testing.T) {
	t.Parallel()

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.PublicInterface(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "public"),
				),
			},
			{
				Config: tmpl.PublicAndVPCInterfaces(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "public"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.1.purpose", "vpc"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image", "interface", "resize_disk", "migration_type"},
			},
			{
				Config: tmpl.VPCAndPublicInterfaces(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "2"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "vpc"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.1.purpose", "public"),
				),
			},
			{
				Config: tmpl.PublicInterface(t, instanceName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "config.0.interface.#", "1"),
					resource.TestCheckResourceAttr(resName, "config.0.interface.0.purpose", "public"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image", "interface", "resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceInstance_migration(t *testing.T) {
	acceptance.LongRunningTest(t)

	t.Parallel()

	rootPass := acctest.RandString(12)

	resName := "linode_instance.foobar"
	var instance linodego.Instance
	instanceName := acctest.RandomWithPrefix("tf_test")

	// Resolve a region to migrate to
	targetRegion, err := acceptance.GetRandomRegionWithCaps(
		[]string{"Linodes"},
		func(v linodego.Region) bool {
			return v.ID != testRegion
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, instanceName, acceptance.PublicKeyMaterial, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
				),
			},
			{
				Config: tmpl.Basic(t, instanceName, acceptance.PublicKeyMaterial, targetRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists(resName, &instance),
					resource.TestCheckResourceAttr(resName, "label", instanceName),
					resource.TestCheckResourceAttr(resName, "type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "region", targetRegion),
				),
			},
			// TODO: Add logic for testing warm migrations once possible
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_pass", "authorized_keys", "image", "resize_disk", "metadata", "migration_type"},
			},
		},
	})
}

func checkInstancePrivateNetworkAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("should have found linode_instance resource %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("should have a Linode ID")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("should have an integer Linode ID: %s", err)
		}

		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
		if err != nil {
			return err
		}

		instanceIPs, err := client.GetInstanceIPAddresses(context.Background(), id)
		if err != nil {
			return err
		}
		if len(instanceIPs.IPv4.Private) == 0 {
			return fmt.Errorf("should have a private ip on Linode ID %d", id)
		}
		return nil
	}
}

type (
	testDiskFunc  func(disk linodego.InstanceDisk) error
	testDisksFunc func(disk []linodego.InstanceDisk) error
)

func testDisk(label string, diskTests ...testDiskFunc) testDisksFunc {
	return func(disks []linodego.InstanceDisk) error {
		for _, disk := range disks {
			if disk.Label == label {
				for _, test := range diskTests {
					if err := test(disk); err != nil {
						return err
					}
				}
				return nil
			}
		}
		return fmt.Errorf("should have found Instance disk with label: %s", label)
	}
}

func testDiskByFS(fs linodego.DiskFilesystem, diskTests ...testDiskFunc) testDisksFunc {
	return func(disks []linodego.InstanceDisk) error {
		for _, disk := range disks {
			if disk.Filesystem == fs {
				for _, test := range diskTests {
					if err := test(disk); err != nil {
						return err
					}
				}
				return nil
			}
		}
		return fmt.Errorf("should have found Instance disk with filesystem: %s", fs)
	}
}

func testDiskExists(diskPtr *linodego.InstanceDisk) testDiskFunc {
	return func(disk linodego.InstanceDisk) error {
		*diskPtr = disk
		return nil
	}
}

func testDiskSize(size int) testDiskFunc {
	return func(disk linodego.InstanceDisk) error {
		if disk.Size != size {
			return fmt.Errorf("should have matching sizes: %d != %d", disk.Size, size)
		}
		return nil
	}
}

func checkInstanceDisks(instance *linodego.Instance, disksTests ...testDisksFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		if instance == nil || instance.ID == 0 {
			return fmt.Errorf("Error fetching disks: invalid Instance argument")
		}

		instanceDisks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)
		if err != nil {
			return fmt.Errorf("Error fetching disks: %s", err)
		}

		if len(instanceDisks) == 0 {
			return fmt.Errorf("No disks")
		}

		for _, tests := range disksTests {
			if err := tests(instanceDisks); err != nil {
				return err
			}
		}

		return nil
	}
}

type (
	testConfigFunc  func(config linodego.InstanceConfig) error
	testConfigsFunc func(config []linodego.InstanceConfig) error
)

// testConfig verifies a labeled config exists and runs many tests against that config
func testConfig(label string, configTests ...testConfigFunc) testConfigsFunc {
	return func(configs []linodego.InstanceConfig) error {
		for _, config := range configs {
			if config.Label == label {
				for _, test := range configTests {
					if err := test(config); err != nil {
						return err
					}
				}
				return nil
			}
		}
		return fmt.Errorf("should have found Instance config with label: %s", label)
	}
}

func testConfigExists(configPtr *linodego.InstanceConfig) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		*configPtr = config
		return nil
	}
}

func testConfigLabel(label string) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if config.Label != label {
			return fmt.Errorf("should have matching labels: %s != %s", config.Label, label)
		}
		return nil
	}
}

func testConfigKernel(kernel string) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if config.Kernel != kernel {
			return fmt.Errorf("should have matching kernels: %s != %s", config.Kernel, kernel)
		}
		return nil
	}
}

func testConfigComments(comments string) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if config.Comments != comments {
			return fmt.Errorf("should have matching comments: %s != %s", config.Comments, comments)
		}
		return nil
	}
}

func testConfigSDADisk(disk *linodego.InstanceDisk) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if disk == nil || config.Devices == nil || config.Devices.SDA == nil || config.Devices.SDA.DiskID != disk.ID {
			return fmt.Errorf("should have SDA with expected disk id")
		}
		return nil
	}
}

func testConfigSDBDisk(disk *linodego.InstanceDisk) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if disk == nil || config.Devices == nil || config.Devices.SDB == nil || config.Devices.SDB.DiskID != disk.ID {
			return fmt.Errorf("should have SDB with expected disk id")
		}
		return nil
	}
}

func testConfigSDBVolume(volume *linodego.Volume) testConfigFunc {
	return func(config linodego.InstanceConfig) error {
		if volume == nil || config.Devices == nil || config.Devices.SDB == nil || config.Devices.SDB.VolumeID != volume.ID {
			return fmt.Errorf("should have SDB with expected volume id")
		}
		return nil
	}
}

func instanceDiskID(disk *linodego.InstanceDisk) string {
	return strconv.Itoa(disk.ID)
}

// checkComputeInstanceConfigs verifies any configs exist and runs config specific tests against a target instance
func checkComputeInstanceConfigs(instance *linodego.Instance, configsTests ...testConfigsFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		if instance == nil || instance.ID == 0 {
			return fmt.Errorf("Error fetching configs: invalid Instance argument")
		}

		instanceConfigs, err := client.ListInstanceConfigs(context.Background(), instance.ID, nil)
		if err != nil {
			return fmt.Errorf("Error fetching configs: %s", err)
		}

		if len(instanceConfigs) == 0 {
			return fmt.Errorf("No configs")
		}

		for _, tests := range configsTests {
			if err := tests(instanceConfigs); err != nil {
				return err
			}
		}

		return nil
	}
}

func checkInstanceDiskExists(instance *linodego.Instance, label string, instanceDisk *linodego.InstanceDisk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		if instance == nil || instance.ID == 0 {
			return fmt.Errorf("Error fetching disks: invalid Instance argument")
		}

		instanceDisks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)
		if err != nil {
			return fmt.Errorf("Error fetching disks: %s", err)
		}

		if len(instanceDisks) == 0 {
			return fmt.Errorf("No disks")
		}

		for _, disk := range instanceDisks {
			if disk.Label == label {
				*instanceDisk = disk
				return nil
			}
		}

		return fmt.Errorf("Disk not found: %s", label)
	}
}

func checkComputeInstanceDisk(instance *linodego.Instance, label string, size int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		if instance == nil || instance.ID == 0 {
			return fmt.Errorf("Error fetching disks: invalid Instance argument")
		}

		instanceDisks, err := client.ListInstanceDisks(context.Background(), instance.ID, nil)
		if err != nil {
			return fmt.Errorf("Error fetching disks: %s", err)
		}

		if len(instanceDisks) == 0 {
			return fmt.Errorf("No disks")
		}

		for _, disk := range instanceDisks {
			if disk.Label == label && disk.Size == size {
				return nil
			}
		}

		return fmt.Errorf("Disk not found: %s", label)
	}
}
