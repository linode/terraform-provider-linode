//go:build integration || producerimagesharegroups

package producerimagesharegroups_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroups/tmpl"
)

func TestAccDataSourceImageShareGroups_basic(t *testing.T) {
	t.Parallel()

	const dsAll = "data.linode_producer_image_share_groups.all"
	const dsByLabel = "data.linode_producer_image_share_groups.by_label"
	const dsByID = "data.linode_producer_image_share_groups.by_id"
	const dsByIsSuspended = "data.linode_producer_image_share_groups.by_is_suspended"

	label1 := acctest.RandomWithPrefix("tf-test")
	label2 := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label1, label2),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(dsAll, "image_share_groups.#", 1),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.id"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.uuid"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.label"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.is_suspended"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.images_count"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.members_count"),
					resource.TestCheckResourceAttrSet(dsAll, "image_share_groups.0.created"),

					resource.TestCheckResourceAttr(dsByLabel, "image_share_groups.#", "1"),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.id"),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.uuid"),
					resource.TestCheckResourceAttr(dsByLabel, "image_share_groups.0.label", label1),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.is_suspended"),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.images_count"),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.members_count"),
					resource.TestCheckResourceAttrSet(dsByLabel, "image_share_groups.0.created"),

					resource.TestCheckResourceAttr(dsByID, "image_share_groups.#", "1"),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.id"),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.uuid"),
					resource.TestCheckResourceAttr(dsByID, "image_share_groups.0.label", label2),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.is_suspended"),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.images_count"),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.members_count"),
					resource.TestCheckResourceAttrSet(dsByID, "image_share_groups.0.created"),

					acceptance.CheckResourceAttrGreaterThan(dsByIsSuspended, "image_share_groups.#", 1),
				),
			},
		},
	})
}
