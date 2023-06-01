package user

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	Username          types.String `tfsdk:"username"`
	SSHKeys           types.List   `tfsdk:"ssh_keys"`
	Email             types.String `tfsdk:"email"`
	Restricted        types.Bool   `tfsdk:"restricted"`
	GlobalGrants      types.List   `tfsdk:"global_grants"`
	DomainGrant       types.Set    `tfsdk:"domain_grant"`
	FirewallGrant     types.Set    `tfsdk:"firewall_grant"`
	ImageGrant        types.Set    `tfsdk:"image_grant"`
	LinodeGrant       types.Set    `tfsdk:"linode_grant"`
	LongviewGrant     types.Set    `tfsdk:"longview_grant"`
	NodebalancerGrant types.Set    `tfsdk:"nodebalancer_grant"`
	StackscriptGrant  types.Set    `tfsdk:"stackscript_grant"`
	VolumeGrant       types.Set    `tfsdk:"volume_grant"`
	DatabaseGrant     types.Set    `tfsdk:"database_grant"`
	ID                types.String `tfsdk:"id"`
}

func (data *DataSourceModel) parseUser(
	ctx context.Context, user *linodego.User,
) diag.Diagnostics {
	data.Username = types.StringValue(user.Username)
	data.Email = types.StringValue(user.Email)
	data.Restricted = types.BoolValue(user.Restricted)

	sshKeys, diags := types.ListValueFrom(ctx, types.StringType, user.SSHKeys)
	if diags.HasError() {
		return diags
	}
	data.SSHKeys = sshKeys

	id, err := json.Marshal(user)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

func (data *DataSourceModel) parseUserGrants(
	ctx context.Context, userGrants *linodego.UserGrants,
) diag.Diagnostics {
	// Domain
	domainGrants, diags := flattenGrantEntities(ctx, userGrants.Domain)
	if diags.HasError() {
		return diags
	}
	data.DomainGrant = *domainGrants

	// Firewall
	firewallGrants, diags := flattenGrantEntities(ctx, userGrants.Firewall)
	if diags.HasError() {
		return diags
	}
	data.FirewallGrant = *firewallGrants

	// Image
	imageGrants, diags := flattenGrantEntities(ctx, userGrants.Image)
	if diags.HasError() {
		return diags
	}
	data.ImageGrant = *imageGrants

	// Linode
	linodeGrants, diags := flattenGrantEntities(ctx, userGrants.Linode)
	if diags.HasError() {
		return diags
	}
	data.LinodeGrant = *linodeGrants

	// Longview
	longviewGrants, diags := flattenGrantEntities(ctx, userGrants.Longview)
	if diags.HasError() {
		return diags
	}
	data.LongviewGrant = *longviewGrants

	// Nodebalancer
	nodebalancerGrants, diags := flattenGrantEntities(ctx, userGrants.NodeBalancer)
	if diags.HasError() {
		return diags
	}
	data.NodebalancerGrant = *nodebalancerGrants

	// Stackscript
	stackscriptGrants, diags := flattenGrantEntities(ctx, userGrants.StackScript)
	if diags.HasError() {
		return diags
	}
	data.StackscriptGrant = *stackscriptGrants

	// Volume
	volumeGrants, diags := flattenGrantEntities(ctx, userGrants.Volume)
	if diags.HasError() {
		return diags
	}
	data.VolumeGrant = *volumeGrants

	// Database
	databaseGrants, diags := flattenGrantEntities(ctx, userGrants.Database)
	if diags.HasError() {
		return diags
	}
	data.DatabaseGrant = *databaseGrants

	// Global

	globalGrants, diags := flattenGlobalGrants(ctx, userGrants.Global)
	if diags.HasError() {
		return diags
	}
	data.GlobalGrants = *globalGrants

	return nil
}

func flattenGlobalGrants(ctx context.Context, grants linodego.GlobalUserGrants) (
	*basetypes.ListValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	if grants.AccountAccess != nil {
		result["account_access"] = types.StringValue(string(*grants.AccountAccess))
	} else {
		result["account_access"] = types.StringValue("")
	}

	result["add_domains"] = types.BoolValue(grants.AddDomains)
	result["add_databases"] = types.BoolValue(grants.AddDatabases)
	result["add_firewalls"] = types.BoolValue(grants.AddFirewalls)
	result["add_images"] = types.BoolValue(grants.AddImages)
	result["add_linodes"] = types.BoolValue(grants.AddLinodes)
	result["add_longview"] = types.BoolValue(grants.AddLongview)
	result["add_nodebalancers"] = types.BoolValue(grants.AddNodeBalancers)
	result["add_stackscripts"] = types.BoolValue(grants.AddStackScripts)
	result["add_volumes"] = types.BoolValue(grants.AddVolumes)
	result["cancel_account"] = types.BoolValue(grants.CancelAccount)
	result["longview_subscription"] = types.BoolValue(grants.LongviewSubscription)

	obj, diag := types.ObjectValue(linodeUserGrantsGlobalObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		linodeUserGrantsGlobalObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}

func flattenGrantEntity(ctx context.Context, entity linodego.GrantedEntity) (
	*basetypes.ObjectValue, diag.Diagnostics,
) {
	result := make(map[string]attr.Value)

	result["id"] = types.Int64Value(int64(entity.ID))
	result["permissions"] = types.StringValue(string(entity.Permissions))

	obj, diag := types.ObjectValue(linodeUserGrantsEntityObjectType.AttrTypes, result)
	if diag.HasError() {
		return nil, diag
	}

	return &obj, nil
}

func flattenGrantEntities(ctx context.Context, entities []linodego.GrantedEntity) (
	*basetypes.SetValue, diag.Diagnostics,
) {
	resultSet := make([]attr.Value, len(entities))

	for i, entity := range entities {
		result, diag := flattenGrantEntity(ctx, entity)
		if diag.HasError() {
			return nil, diag
		}

		resultSet[i] = result
	}

	result, diag := basetypes.NewSetValue(
		linodeUserGrantsEntityObjectType,
		resultSet,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &result, nil
}
