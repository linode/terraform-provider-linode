//go:build integration || lkeversions

package lkeversions_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversions/tmpl"
)

func TestAccDataSourceLinodeLkeVersions_NoTier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_versions.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataNoTier(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("versions").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("versions").AtSliceIndex(0).AtMapKey("tier"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("tier"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

func TestAccDataSourceLinodeLkeVersions_Tier(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_lke_versions.foobar"
	tier := "enterprise"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataTier(t, tier),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("versions").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("versions").AtSliceIndex(0).AtMapKey("tier"),
						knownvalue.StringExact(tier),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("tier"),
						knownvalue.StringExact(tier),
					),
				},
			},
		},
	})
}
