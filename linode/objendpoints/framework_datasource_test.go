//go:build integration || objendpoints

package objendpoints_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/objendpoints/tmpl"
)

func TestAccDataSourceObjectStorageEndpoints_basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_object_storage_endpoints.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("endpoints"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("endpoints").AtSliceIndex(0).AtMapKey("endpoint_type"),
						knownvalue.StringRegexp(regexp.MustCompile(`^E\d$`)),
					),
				},
			},
		},
	})
}

func TestAccDataSourceObjectStorageEndpoints_filter(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_object_storage_endpoints.test"
	targetEndpointType := "E2"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilter(t, targetEndpointType),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("endpoints").AtSliceIndex(0).AtMapKey("endpoint_type"),
						knownvalue.StringExact(targetEndpointType),
					),
				},
			},
		},
	})
}
