//go:build integration || nfsspaces

package nfsspaces_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/nfsspaces/tmpl"
)

func TestAccDataSourceNFSSpaces_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nfs_spaces.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, label),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_spaces").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label)),
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_spaces").AtSliceIndex(0).AtMapKey("status"),
						knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccDataSourceNFSSpaces_filterLabel(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_nfs_spaces.foobar"
	label := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterLabel(t, label),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(dataSourceName,
						tfjsonpath.New("nfs_spaces").AtSliceIndex(0).AtMapKey("label"),
						knownvalue.StringExact(label)),
				},
			},
		},
	})
}
