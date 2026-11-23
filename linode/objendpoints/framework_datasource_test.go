//go:build integration || objendpoints

package objendpoints_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/objendpoints/tmpl"
)

func stringExactOrNull(s string) knownvalue.Check {
	if s == "" {
		return knownvalue.Null()
	}
	return knownvalue.StringExact(s)
}

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
						knownvalue.StringRegexp(regexp.MustCompile(`^(E|IA)\d$`)),
					),
				},
			},
		},
	})
}

func TestAccDataSourceObjectStorageEndpoints_ia(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.linode_object_storage_endpoints.test"
	endpoint := acceptance.GetRandomObjectStorageEndpointOrSkip(t, linodego.ObjectStorageEndpointIA0)
	// The API leaves out the IA fields while the account's IA API flag is off.
	if len(endpoint.SupportedKeyClasses) == 0 {
		t.Skip("skipping: the IA0 endpoint lists no supported_key_classes")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataIA(t, endpoint.Region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("endpoints"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"endpoint_type": knownvalue.StringExact(string(linodego.ObjectStorageEndpointIA0)),
								"region":        knownvalue.StringExact(endpoint.Region),
								"s3_endpoint":   knownvalue.StringExact(*endpoint.S3Endpoint),
								"supported_key_classes": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact(string(linodego.ObjectStorageKeyClassIA)),
								}),
								"country":      stringExactOrNull(endpoint.Country),
								"region_label": stringExactOrNull(endpoint.RegionLabel),
							}),
						}),
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
