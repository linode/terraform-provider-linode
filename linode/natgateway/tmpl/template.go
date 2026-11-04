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
	UseAltSubnet     bool
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

// Complex renders a manual (autoscaling false) NAT gateway with one reserved address, attached to a VPC subnet.
func Complex(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_complex", TemplateData{
		Label:  label,
		Region: region,
	})
}

// ComplexAltSubnet renders the complex config with the gateway attached to the alt subnet instead.
func ComplexAltSubnet(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_complex", TemplateData{
		Label:        label,
		Region:       region,
		UseAltSubnet: true,
	})
}

// DataBasic renders the complex NAT gateway config and a data source that reads it.
func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_data_basic", TemplateData{
		Label:  label,
		Region: region,
	})
}
