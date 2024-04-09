//go:build integration || volumes

package volumes_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/volumes/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_volumes", &resource.Sweeper{
		Name: "linode_volumes",
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

func TestAccDataSourceVolume_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_volumes.volume"
	volumeName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, volumeName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "volumes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "volumes.0.region", testRegion),
				),
			},
		},
	})
}
