package rdns

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

type ResourceModel struct {
	Address          types.String `tfsdk:"address"`
	RDNS             types.String `tfsdk:"rdns"`
	WaitForAvailable types.Bool   `tfsdk:"wait_for_available"`
	ID               types.String `tfsdk:"id"`
}

func (rm *ResourceModel) parseConfiguredAttributes(ip *linodego.InstanceIP) {
	rm.Address = types.StringValue(ip.Address)

	if !rm.RDNS.Equal(types.StringValue(ip.RDNS)) {
		rm.RDNS = types.StringValue(ip.RDNS)
	}
}

func (rm *ResourceModel) parseComputedAttributes(ip *linodego.InstanceIP) {
	rm.ID = types.StringValue(ip.Address)
}
