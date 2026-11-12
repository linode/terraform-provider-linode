//go:build integration || natgateway

package natgateway_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/natgateway/tmpl"
)

const (
	natAutoTrue     = "linode_natgateway.explicit_autoscaling_true"
	natAutoFalse    = "linode_natgateway.explicit_autoscaling_false"
	natDefaultEmpty = "linode_natgateway.default_empty_addresses"
	natDefaultPop   = "linode_natgateway.default_populated_addresses"
	natComplex      = "linode_natgateway.complex"
)

// TestAccResourceNATGateway_update creates the four NAT Gateway variants, changes their
// labels, and swaps the manual gateway's address. Each update must happen in place.
func TestAccResourceNATGateway_update(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The firewall label is derived from this value and must begin with a letter.
	label := "tf-" + acctest.RandString(10)
	updatedLabel := "tf-" + acctest.RandString(10)

	// natGatewayStateChecks returns the expected state for each variant. base is the
	// label prefix used in the template for this step.
	natGatewayStateChecks := func(base string) []statecheck.StateCheck {
		return []statecheck.StateCheck{
			// Explicit autoscaling: the configured value is kept in state.
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("use_autoscaling"), knownvalue.Bool(true)),
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-autoscaling-true")),
			statecheck.ExpectKnownValue(natAutoTrue,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(0)),

			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("use_autoscaling"), knownvalue.Bool(false)),
			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-autoscaling-false")),
			statecheck.ExpectKnownValue(natAutoFalse,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),

			// Implied autoscaling: not configured, so not stored (the API does not return it).
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("use_autoscaling"), knownvalue.Null()),
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-default-empty")),
			statecheck.ExpectKnownValue(natDefaultEmpty,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(0)),

			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("use_autoscaling"), knownvalue.Null()),
			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("label"), knownvalue.StringExact(base+"-default-populated")),
			statecheck.ExpectKnownValue(natDefaultPop,
				tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create all four variants.
				Config:            tmpl.Basic(t, label, region),
				ConfigStateChecks: natGatewayStateChecks(label),
			},
			{
				// Update labels only. Every variant must update in place (no replacement),
				// and use_autoscaling / addresses must be unchanged.
				Config: tmpl.Basic(t, updatedLabel, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(natAutoTrue, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natAutoFalse, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natDefaultEmpty, plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction(natDefaultPop, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: natGatewayStateChecks(updatedLabel),
			},
			{
				// Swap the manual gateway's address from address 1 to address 3. Address 3 is
				// not used by any other gateway. Must be an in-place update.
				Config: tmpl.AddressSwap(t, updatedLabel, region),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(natAutoFalse, plancheck.ResourceActionUpdate),
					},
				},
				ConfigStateChecks: natGatewayStateChecks(updatedLabel),
			},
			{
				// Assign an address to the autoscaling-true gateway. The provider must
				// reject this during validation, before any API call is made.
				Config:      tmpl.AutoscaleAddressUpdate(t, updatedLabel, region),
				ExpectError: regexp.MustCompile(`(?s)cannot be configured when\s+autoscaling is enabled`),
			},
			{
				// Restore the valid configuration. The post-test destroy reuses the
				// configuration left by the final step, which must therefore be valid.
				Config:            tmpl.AddressSwap(t, updatedLabel, region),
				ConfigStateChecks: natGatewayStateChecks(updatedLabel),
			},
		},
	})
}

// TestAccResourceNATGateway_autoscalingEmptyAddresses asserts that an explicitly configured
// empty addresses set is rejected for an autoscaling gateway. Autoscaling allocates addresses
// that Read stores in the attribute, so such a configuration could never be reconciled.
//
// Both configurations fail validation, so no resources are created.
func TestAccResourceNATGateway_autoscalingEmptyAddresses(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityNATGateway}, "core")
	if err != nil {
		t.Fatal(err)
	}

	label := "tf-" + acctest.RandString(10)
	expectedError := regexp.MustCompile(`(?s)cannot be configured when\s+autoscaling is enabled`)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// use_autoscaling is explicitly true.
				Config:      tmpl.AutoscalingEmptyAddresses(t, label, region),
				ExpectError: expectedError,
			},
			{
				// use_autoscaling is omitted, so an empty addresses set implies autoscaling.
				Config:      tmpl.ImpliedAutoscalingEmptyAddresses(t, label, region),
				ExpectError: expectedError,
			},
			{
				// The post-test destroy reuses the configuration left by the final step.
				// An empty configuration keeps that destroy valid, since the preceding
				// steps never created anything.
				Config: "# No resources. See the comment on this test step.\n",
			},
		},
	})
}

// TestAccResourceNATGateway_complex creates a manual (autoscaling false) NAT Gateway with one
// reserved address.
func TestAccResourceNATGateway_complex(t *testing.T) {
	t.Parallel()

	region, err := acceptance.GetRandomRegionWithCaps(
		[]linodego.RegionCapability{linodego.CapabilityVPCs, linodego.CapabilityNATGateway}, "core")
	if err != nil {
		t.Fatal(err)
	}

	// The label must begin with a letter.
	label := "tf-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Complex(t, label, region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("label"), knownvalue.StringExact(label+"-complex")),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("use_autoscaling"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("addresses"), knownvalue.SetSizeExact(1)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("vpc_subnet"), knownvalue.Null()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("region"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("default_ports_per_interface"), knownvalue.Int64Exact(4096)),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("address_autoscale_max"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("portset_assignments"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("portset_capacity"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("created"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(natComplex,
						tfjsonpath.New("updated"), knownvalue.NotNull()),
				},
			},
		},
	})
}
