package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	UseSwapAddress   bool
	AutoscaleAddress bool
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_basic", TemplateData{
		Label:  label,
		Region: region,
	})
}

// AddressSwap renders the config where the manual gateway uses address 3 instead of address 1.
func AddressSwap(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_basic", TemplateData{
		Label:          label,
		Region:         region,
		UseSwapAddress: true,
	})
}

// AutoscaleAddressUpdate renders the swap config and also assigns address 1 to the autoscaling-true
// gateway. The API must reject this.
func AutoscaleAddressUpdate(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_basic", TemplateData{
		Label:            label,
		Region:           region,
		UseSwapAddress:   true,
		AutoscaleAddress: true,
	})
}

// AutoscalingEmptyAddresses renders an autoscaling gateway that explicitly declares an empty
// addresses set. The provider must reject this during validation.
func AutoscalingEmptyAddresses(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_autoscaling_empty_addresses", TemplateData{
		Label:  label,
		Region: region,
	})
}

// ImpliedAutoscalingEmptyAddresses renders a gateway that omits use_autoscaling and declares an
// empty addresses set, which implies autoscaling. The provider must reject this during validation.
func ImpliedAutoscalingEmptyAddresses(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_implied_autoscaling_empty_addresses", TemplateData{
		Label:  label,
		Region: region,
	})
}

// Complex renders a manual (autoscaling false) NAT gateway with one reserved address, attached to a VPC subnet.
func Complex(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_complex", TemplateData{
		Label:  label,
		Region: region,
	})
}

// DataBasic renders the complex NAT gateway config and a data source that reads it.
func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_data_basic", TemplateData{
		Label:  label,
		Region: region,
	})
}
