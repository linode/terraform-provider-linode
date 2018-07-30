package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestDetectVolumeIDChange(t *testing.T) {
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

func TestAccLinodeVolumeBasic(t *testing.T) {
	t.Parallel()

	resName := "linode_volume.foobar"
	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr(resName, "label", volumeName),
				),
			},

			resource.TestStep{
				ResourceName: resName,
				ImportState:  true,
			},
		},
	})
}

func TestAccLinodeVolumeUpdate(t *testing.T) {
	t.Parallel()

	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigUpdates(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", fmt.Sprintf("%s_renamed", volumeName)),
				),
			},
		},
	})
}

func TestAccLinodeVolumeResized(t *testing.T) {
	t.Parallel()

	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigResized(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "size", "30"),
				),
			},
		},
	})
}

func TestAccLinodeVolumeAttached(t *testing.T) {
	t.Parallel()

	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckNoResourceAttr("linode_volume.foobar", "linode_id"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttrSet("linode_instance.foobar", "id"),
					resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
				),
			},
		},
	})
}

func TestAccLinodeVolumeDetached(t *testing.T) {
	t.Parallel()

	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrPair("linode_instance.foobar", "id", "linode_volume.foobar", "linode_id"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
				),
			},
		},
	})
}

func TestAccLinodeVolumeReattachedBetweenInstances(t *testing.T) {
	t.Parallel()

	var volumeName = fmt.Sprintf("tf_test_%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigAttached(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			resource.TestStep{
				Config: testAccCheckLinodeVolumeConfigReattachedBetweenInstances(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeVolumeExists,
					resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobaz", "id"),
				),
			},
		},
	})
}

func testAccCheckLinodeVolumeExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_volume" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)

		_, err = client.GetVolume(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Volume %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeVolumeDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Failed to get Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_volume" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Failed parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetVolume(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Volume with id %d still exists", id)
		}

		if apiErr, ok := err.(linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Failed to request Linode Volume with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeVolumeConfigBasic(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
}`, volume)
}

func testAccCheckLinodeVolumeConfigUpdates(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s_renamed"
	region = "us-west"
}`, volume)
}

func testAccCheckLinodeVolumeConfigResized(volume string) string {
	return fmt.Sprintf(`
resource "linode_volume" "foobar" {
	label = "%s_renamed"
	region = "us-west"
	size = 30
}`, volume)
}

func testAccCheckLinodeVolumeConfigAttached(volume string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	root_password = "%s"
	region = "us-west"
}
	
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	linode_id = "${linode_instance.foobar.id}"
}`, volume, volume)
}

func testAccCheckLinodeVolumeConfigReattachedBetweenInstances(volume string) string {
	return fmt.Sprintf(`
resource "linode_instance" "foobar" {
	root_password = "%s"
	region = "us-west"
}

resource "linode_instance" "foobaz" {
	root_password = "%s"
	region = "us-west"
}
	
resource "linode_volume" "foobar" {
	label = "%s"
	region = "us-west"
	linode_id = "${linode_instance.foobaz.id}"
}`, volume, volume, volume)
}
