//go:build (integration && long_running) || instance

package instance_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instance/tmpl"
	"testing"
)

var testRegion string

func TestAccResourceInstance_migration(t *testing.T) {
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
