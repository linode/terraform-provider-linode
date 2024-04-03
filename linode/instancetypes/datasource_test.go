//go:build integration || instancetypes

package instancetypes_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetypes/tmpl"
)

func TestAccDataSourceInstanceTypes_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "types.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "types.0.id", "g6-standard-2"),
					resource.TestCheckResourceAttr(resourceName, "types.0.label", "Linode 4GB"),
					resource.TestCheckResourceAttr(resourceName, "types.0.class", "standard"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.disk"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.network_out"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.memory"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.transfer"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.price.0.monthly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.addons.0.backups.0.price.0.hourly"),
					resource.TestCheckResourceAttrSet(resourceName, "types.0.addons.0.backups.0.price.0.monthly"),
				),
			},
		},
	})
}

func TestAccDataSourceInstanceTypes_substring(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataSubstring(t),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "types.#", 1),
					acceptance.CheckResourceAttrContains(resourceName, "types.0.label", "Linode"),
				),
			},
		},
	})
}

func TestAccDataSourceInstanceTypes_regex(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataRegex(t),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "types.#", 1),
					acceptance.CheckResourceAttrContains(resourceName, "types.0.label", "Dedicated"),
				),
			},
		},
	})
}

func TestAccDataSourceInstanceTypes_byClass(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_instance_types.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataByClass(t),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(resourceName, "types.#", 0),
					acceptance.CheckResourceAttrContains(resourceName, "types.0.label", "Linode"),
				),
			},
		},
	})
}
