//go:build integration || volume

package volume_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/volume/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_volume", &resource.Sweeper{
		Name: "linode_volume",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{"Block Storage"})
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

func TestAccResourceVolume_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_volume.foobar"
	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttrSet(resName, "size"),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
					resource.TestCheckResourceAttr(resName, "region", testRegion),
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

	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}
	resName := "linode_volume.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", volumeName),
				),
			},
			{
				Config: tmpl.Updates(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_r", volumeName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
			{
				// Expect no changes on tags case changed
				Config:   tmpl.UpdatesTagsCaseChange(t, volumeName, testRegion),
				PlanOnly: true,
			},
		},
	})
}

func TestAccResourceVolume_resized(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config: tmpl.Resized(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "size", "30"),
					resource.TestCheckResourceAttr("linode_volume.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceVolume_attached(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
				),
			},
			{
				Config: tmpl.Attached(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttrSet("linode_instance.foobar", "id"),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				ResourceName:            "linode_volume.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id"},
				Check:                   resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
			},
		},
	})
}

func TestAccResourceVolume_detached(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Attached(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
				),
			},
			{
				Config:                  tmpl.Attached(t, volumeName, testRegion),
				ResourceName:            "linode_volume.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"linode_id"},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("linode_volume.foobar", "linode_id", "0"),
					resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobar", "id"),
				),
			},
		},
	})
}

func TestAccResourceVolume_reattachedBetweenInstances(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	volume := linodego.Volume{}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Attached(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				Config: tmpl.ReAttached(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
				),
			},
			{
				ResourceName:            "linode_instance.foobar",
				Check:                   resource.TestCheckResourceAttrPair("linode_volume.foobaz", "linode_id", "linode_instance.foobar", "id"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
			{
				ResourceName:            "linode_instance.foobaz",
				Check:                   resource.TestCheckResourceAttrPair("linode_volume.foobar", "linode_id", "linode_instance.foobaz", "id"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resize_disk", "migration_type"},
			},
		},
	})
}

func TestAccResourceVolume_cloned(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")

	var instance linodego.Instance
	var instance2 linodego.Instance

	var volume linodego.Volume
	var volume2 linodego.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ClonedStep1(t, volumeName, testRegion, acceptance.PublicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					acceptance.CheckInstanceExists("linode_instance.foobar2", &instance2),

					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),
				),
			},
			{
				Config: tmpl.ClonedStep1(t, volumeName, testRegion, acceptance.PublicKeyMaterial),
				PreConfig: func() {
					outBuffer := new(bytes.Buffer)

					client := acceptance.GetSSHClient(t, "root", instance.IPv4[0].String())

					defer client.Close()
					session, err := client.NewSession()
					if err != nil {
						t.Fatalf("failed to establish SSH session: %s", err)
					}

					session.Stdout = outBuffer

					// Format the first volume and drop a file onto it
					err = session.Run(fmt.Sprintf(scriptFormatDrive,
						volume.FilesystemPath, volume.FilesystemPath, volume.FilesystemPath))
					if err != nil {
						t.Fatalf("failed to format and mount volume: %s", err)
					}
				},
			},
			{
				// Clone the volume
				Config: tmpl.ClonedStep2(t, volumeName, testRegion, acceptance.PublicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckInstanceExists("linode_instance.foobar", &instance),
					acceptance.CheckInstanceExists("linode_instance.foobar2", &instance2),

					acceptance.CheckVolumeExists("linode_volume.foobar", &volume),
					resource.TestCheckResourceAttr("linode_volume.foobar", "label", volumeName),
					resource.TestCheckResourceAttrSet("linode_volume.foobar", "linode_id"),

					acceptance.CheckVolumeExists("linode_volume.foobar-cloned", &volume2),
					resource.TestCheckResourceAttr("linode_volume.foobar-cloned", "label", volumeName+"-c"),
					resource.TestCheckResourceAttrSet("linode_volume.foobar-cloned", "linode_id"),
				),
			},
			{
				Config: tmpl.ClonedStep2(t, volumeName, testRegion, acceptance.PublicKeyMaterial),
				PreConfig: func() {
					outBuffer := new(bytes.Buffer)
					client := acceptance.GetSSHClient(t, "root", instance2.IPv4[0].String())

					defer client.Close()
					session, err := client.NewSession()
					if err != nil {
						t.Fatalf("failed to establish SSH session: %s", err)
					}

					session.Stdout = outBuffer

					// Check that the file was cloned onto the new volume
					err = session.Run(fmt.Sprintf(scriptCheckCloneExists,
						volume2.FilesystemPath, volume2.FilesystemPath))
					if err != nil {
						t.Fatalf("failed to check for cloned file: %s", err)
					}
				},
			},
		},
	})
}

const scriptFormatDrive = `
until [ -e "%s" ]; do sleep .1; done && \
mkfs.ext4 "%s" && \
mkdir -p /mnt/vol && \
mount "%s" "/mnt/vol" && \
touch /mnt/vol/itworks.txt && \
umount /mnt/vol
`

const scriptCheckCloneExists = `
until [ -e "%s" ]; do sleep .1; done && \
echo $? && \
mkdir -p /mnt/vol && \
echo $? && \
mount "%s" "/mnt/vol" && \
echo $? && \
test -f /mnt/vol/itworks.txt && \
echo $? && \
umount /mnt/vol
`
