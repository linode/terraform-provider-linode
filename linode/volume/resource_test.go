package volume_test

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
	"github.com/linode/terraform-provider-linode/linode/volume"
)

func init() {
	resource.AddTestSweepers("linode_volume", &resource.Sweeper{
		Name: "linode_volume",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	volumes, err := client.ListVolumes(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting volumes: %s", err)
	}
	for _, volume := range volumes {
		if !acceptance.ShouldSweep(prefix, volume.Label) {
			continue
		}
		err := client.DeleteVolume(context.Background(), volume.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", volume.Label, err)
		}
	}

	return nil
}

func TestDetectVolumeIDChange(t *testing.T) {
	t.Parallel()
	var have, want *int
	var one, two *int
	oneValue, twoValue := 1, 2
	one, two = &oneValue, &twoValue

	if have, want = nil, nil; volume.DetectVolumeIDChange(have, want) {
		t.Errorf("should not detect change when both are nil")
	}
	if have, want = nil, one; !volume.DetectVolumeIDChange(have, want) {
		t.Errorf("should detect change when have is nil and want is not nil")
	}
	if have, want = one, nil; !volume.DetectVolumeIDChange(have, want) {
		t.Errorf("should detect change when want is nil and have is not nil")
	}
	if have, want = one, two; !volume.DetectVolumeIDChange(have, want) {
		t.Errorf("should detect change when values differ")
	}
}

func TestAccResourceVolume_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_volume.foobar"
	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
					resource.TestCheckResourceAttr(resName, "region", "us-west"),
					resource.TestCheckResourceAttr(resName, "linode_id", "0"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
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

func TestAccResourceVolume_update(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}
	var resName = "linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
				),
			},
			{
				Config: resourceConfigUpdates(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", volumeName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
		},
	})
}

func TestAccResourceVolume_resized(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config: resourceConfigVolumeResized(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "size", "30"),
					resource.TestCheckResourceAttr("linode_volume.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceVolume_attached(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
				),
			},
			{
				Config: resourceConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttrSet("linode_instance.foobar", "id"),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				ResourceName:      "linode_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				Check:             resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
			},
		},
	})
}

func TestAccResourceVolume_detached(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config:            resourceConfigAttached(volumeName),
				ResourceName:      "linode_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				Check:             resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
			},
			{
				Config:            resourceConfigBasic(volumeName),
				ResourceName:      "linode_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				Check:             resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
			},
		},
	})
}

func TestAccResourceVolume_reattachedBetweenInstances(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				Config: resourceConfigReattachedBetweenInstances(volumeName),
				Check: resource.ComposeTestCheckFunc(
					checkVolumeExists("linode_volume.foobar", &volume),
				),
			},
			{
				ResourceName:      "linode_instance.foobar",
				Check:             resource.TestCheckResourceAttrPair("linode_volume.foobaz", "linode_id", "linode_instance.foobar", "id"),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "linode_instance.foobaz",
				Check:             resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobaz", "id"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkVolumeExists(name string, volume *linodego.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetVolume(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Volume %s: %s", rs.Primary.Attributes["label"], err)
		}

		*volume = *found

		return nil
	}
}

func checkVolumeDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_volume" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetVolume(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Volume with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Volume with id %d", id)
		}
	}

	return nil
}

func resourceConfigBasic(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	tags = ["tf_test"]
}`, volume)
}

func resourceConfigUpdates(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s_r"
	region = "us-west"
	tags = ["tf_test", "tf_test_2"]
}`, volume)
}

func resourceConfigVolumeResized(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	size = 30
}`, volume)
}

func resourceConfigAttached(volume string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	type = "g6-nanode-1"
	region = "us-west"

	config {
		label = "config"
		kernel = "linode/latest-64bit"
		devices {
			sda {
				volume_id = "${linode_volume.foobar.id}"
			}
		}
	}
}

resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
}`, volume)
}

func resourceConfigReattachedBetweenInstances(volume string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	type = "g6-nanode-1"
	region = "us-west"

	config {
		label = "config"
		kernel = "linode/latest-64bit"
		devices {
			sda {
				volume_id = "${linode_volume.foobaz.id}"
			}
		}
	}
}

resource "linode_instance" "foobaz" {
	type = "g6-nanode-1"
	region = "us-west"

	config {
		label = "config"
		kernel = "linode/latest-64bit"
		devices {
			sda {
				volume_id = "${linode_volume.foobar.id}"
			}
		}
	}
}

resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
}

resource "linode_volume" "foobaz" {
	label = "%s_baz"
	region = "us-west"
}
`, volume, volume)
}
