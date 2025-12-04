//go:build integration || volume

package volume_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/volume/tmpl"
)

// Default encryption (omitted) should be enabled (provider derives default at create-time)
func TestAccResourceVolume_defaultEncryptionEnabled_Derived(t *testing.T) {
	 t.Parallel()

	 volumeName := acctest.RandomWithPrefix("tf_test")
	 resName := "linode_volume.foobar"

	 // Choose a random core region without checking capabilities
	 targetRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	 if err != nil {
		 t.Fatal(err)
	 }

	volume := linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				// Basic template omits encryption
				Config: tmpl.Basic(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "region", targetRegion),
					resource.TestCheckResourceAttr(resName, "encryption", "enabled"),
				),
			},
		},
	})
}

// Explicit encryption enabled (resource test)
func TestAccResourceVolume_encryptionExplicitEnabled(t *testing.T) {
	 t.Parallel()

	 volumeName := acctest.RandomWithPrefix("tf_test")
	 resName := "linode_volume.foobar"

	 // Choose a random core region without checking capabilities
	 targetRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	 if err != nil {
		 t.Fatal(err)
	 }

	volume := linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithBlockStorageEncryption(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "encryption", "enabled"),
				),
			},
		},
	})
}

// Explicit encryption disabled (resource test)
func TestAccResourceVolume_encryptionExplicitDisabled(t *testing.T) {
	 t.Parallel()

	 volumeName := acctest.RandomWithPrefix("tf_test")
	 resName := "linode_volume.foobar"

	 // Choose a random core region without checking capabilities
	 targetRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	 if err != nil {
		 t.Fatal(err)
	 }

	volume := linodego.Volume{}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithBlockStorageEncryptionDisabled(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &volume),
					resource.TestCheckResourceAttr(resName, "encryption", "disabled"),
				),
			},
		},
	})
}

// Changing encryption forces replacement (verify ID changes)
func TestAccResourceVolume_encryptionChangeForcesReplace(t *testing.T) {
	t.Parallel()

	volumeName := acctest.RandomWithPrefix("tf_test")
	resName := "linode_volume.foobar"

targetRegion, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		t.Fatal(err)
	}

	var v linodego.Volume
	var firstID int

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             acceptance.CheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithBlockStorageEncryptionDisabled(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &v),
					resource.TestCheckResourceAttr(resName, "encryption", "disabled"),
					func(_ *terraform.State) error { firstID = v.ID; return nil },
				),
			},
			{
				Config: tmpl.DataWithBlockStorageEncryption(t, volumeName, targetRegion),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckVolumeExists(resName, &v),
					resource.TestCheckResourceAttr(resName, "encryption", "enabled"),
					func(_ *terraform.State) error {
						if v.ID == firstID {
							return fmt.Errorf("expected replacement, ID unchanged: %d", v.ID)
						}
						return nil
					},
				),
			},
		},
	})
}
