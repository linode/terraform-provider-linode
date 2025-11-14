//go:build integration || producerimagesharegroupimageshares

package producerimagesharegroupimageshares_test

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupimageshares/tmpl"
)

func TestAccDataSourceImageShareGroupImageShares_basic(t *testing.T) {
	t.Parallel()

	const dsAll = "data.linode_producer_image_share_group_image_shares.all"
	const dsByID = "data.linode_producer_image_share_group_image_shares.by_id"
	const dsByLabel = "data.linode_producer_image_share_group_image_shares.by_label"

	label := acctest.RandomWithPrefix("tf_test")
	instanceLabel := acctest.RandomWithPrefix("tf_test")

	instanceRegion, err := acceptance.GetRandomRegionWithCaps([]string{}, "core")
	if err != nil {
		log.Fatal(err)
	}

	imageLabel1 := acctest.RandomWithPrefix("tf-test")
	imageLabel2 := acctest.RandomWithPrefix("tf-test")
	shareGroupLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label, instanceLabel, instanceRegion, imageLabel1, imageLabel2, shareGroupLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsAll, "image_shares.#", "2"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.0.label", "image_one_label"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.0.description", "image one description"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.1.label", "image_two_label"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.1.description", "image two description"),

					resource.TestCheckResourceAttr(dsByID, "image_shares.#", "1"),
					resource.TestCheckResourceAttr(dsByID, "image_shares.0.label", "image_one_label"),
					resource.TestCheckResourceAttr(dsByID, "image_shares.0.description", "image one description"),

					resource.TestCheckResourceAttr(dsByLabel, "image_shares.#", "1"),
					resource.TestCheckResourceAttr(dsByLabel, "image_shares.0.label", "image_two_label"),
					resource.TestCheckResourceAttr(dsByLabel, "image_shares.0.description", "image two description"),
				),
			},
		},
	})
}
