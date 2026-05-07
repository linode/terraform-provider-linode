//go:build integration || iamentities

package iamentities_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/iamentities/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_volume", &resource.Sweeper{
		Name: "linode_volume",
		F:    volumeSweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps([]string{linodego.CapabilityBlockStorage}, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func volumeSweep(prefix string) error {
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

func TestAccDataSourceIAMEntities_basic(t *testing.T) {
	t.Parallel()

	// IAM Tests need to be opted into, iam accounts do not support all existing user endpoints as they will be replacing some of them
	acceptance.OptInTest(t)

	label := acctest.RandomWithPrefix("tf_test")
	resName := "data.linode_iam_entities.test_iam_entities"
	// entities doesn't support filtering at the moment so i'm only
	// checking type incase test account has multiple volumes
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, testRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "entities.0.id"),
					resource.TestCheckResourceAttr(resName, "entities.0.type", "volume"),
				),
			},
		},
	})
}
