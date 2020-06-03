package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_volume", &resource.Sweeper{
		Name: "linode_volume",
		F:    testSweepLinodeVolume,
	})
}

func testSweepLinodeVolume(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	volumes, err := client.ListVolumes(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting volumes: %s", err)
	}
	for _, volume := range volumes {
		if !shouldSweepAcceptanceTestResource(prefix, volume.Label) {
			continue
		}
		err := client.DeleteVolume(context.Background(), volume.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", volume.Label, err)
		}
	}

	return nil
}

func TestAccLinodeVolume_detectVolumeIDChange(t *testing.T) {
	t.Parallel()
	var have, want *int
	var one, two *int
	oneValue, twoValue := 1, 2
	one, two = &oneValue, &twoValue

	if have, want = nil, nil; detectVolumeIDChange(have, want) {
		t.Errorf("should not detect change when both are nil")
	}
	if have, want = nil, one; !detectVolumeIDChange(have, want) {
		t.Errorf("should detect change when have is nil and want is not nil")
	}
	if have, want = one, nil; !detectVolumeIDChange(have, want) {
		t.Errorf("should detect change when want is nil and have is not nil")
	}
	if have, want = one, two; !detectVolumeIDChange(have, want) {
		t.Errorf("should detect change when values differ")
	}
}

func TestAccLinodeVolume_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_volume.foobar"
	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
					resource.TestCheckResourceAttr(resName, "region", "us-west"),
					resource.TestCheckResourceAttr(resName, "linode_id", "0"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
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

func TestAccLinodeVolume_update(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}
	var resName = "linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
				),
			},
			{
				Config: testAccCheckLinodeVolumeConfigUpdates(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", volumeName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.4106436895", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.2667398925", "tf_test_2"),
				),
			},
		},
	})
}

func TestAccLinodeVolume_resized(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config: testAccCheckLinodeVolumeConfigResized(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "size", "30"),
					resource.TestCheckResourceAttr("linode_volume.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccLinodeVolume_attached(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
				),
			},
			{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
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

func TestAccLinodeVolume_detached(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config:            testAccCheckLinodeVolumeConfigAttached(volumeName),
				ResourceName:      "linode_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				Check:             resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
			},
			{
				Config:            testAccCheckLinodeVolumeConfigBasic(volumeName),
				ResourceName:      "linode_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				Check:             resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
			},
		},
	})
}

func TestAccLinodeVolume_reattachedBetweenInstances(t *testing.T) {
	t.Parallel()

	var volumeName = acctest.RandomWithPrefix("tf_test")
	var volume = linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				Config: testAccCheckLinodeVolumeConfigReattachedBetweenInstances(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists("linode_volume.foobar", &volume),
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

func testAccCheckLinodeVolumeExists(name string, volume *linodego.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(linodego.Client)

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

func testAccCheckLinodeVolumeDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
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

func testAccCheckLinodeVolumeConfigBasic(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	tags = ["tf_test"]
}`, volume)
}

func testAccCheckLinodeVolumeConfigUpdates(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s_r"
	region = "us-west"
	tags = ["tf_test", "tf_test_2"]
}`, volume)
}

func testAccCheckLinodeVolumeConfigResized(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	size = 30
}`, volume)
}

func testAccCheckLinodeVolumeConfigAttached(volume string) string {
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

func testAccCheckLinodeVolumeConfigReattachedBetweenInstances(volume string) string {
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
