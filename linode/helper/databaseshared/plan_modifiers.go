package databaseshared

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/linode/terraform-provider-linode/v3/linode/helper/stringplanmodifiers"
)

// HostStringPlanModifiers is the single source of truth for DB host plan modifiers.
//
// It prevents computed host fields (e.g. host_primary/host_secondary/host_standby)
// from being replaced with unknown during planning unless the VPC/private network
// configuration or instance type changes.
var HostStringPlanModifiers = []planmodifier.String{
	stringplanmodifiers.UseStateForUnknownUnlessTheseChanged(
		path.MatchRoot("private_network"),
		path.MatchRoot("type"),
		path.MatchRoot("cluster_size"),
	),
}
