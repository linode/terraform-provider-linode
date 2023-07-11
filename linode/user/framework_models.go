package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/linode/linodego"
)

type DataSourceModel struct {
	Username            types.String `tfsdk:"username"`
	SSHKeys             types.List   `tfsdk:"ssh_keys"`
	Email               types.String `tfsdk:"email"`
	Restricted          types.Bool   `tfsdk:"restricted"`
	GlobalGrants        types.List   `tfsdk:"global_grants"`
	DomainGrant         types.Set    `tfsdk:"domain_grant"`
	FirewallGrant       types.Set    `tfsdk:"firewall_grant"`
	ImageGrant          types.Set    `tfsdk:"image_grant"`
	LinodeGrant         types.Set    `tfsdk:"linode_grant"`
	LongviewGrant       types.Set    `tfsdk:"longview_grant"`
	NodebalancerGrant   types.Set    `tfsdk:"nodebalancer_grant"`
	StackscriptGrant    types.Set    `tfsdk:"stackscript_grant"`
	VolumeGrant         types.Set    `tfsdk:"volume_grant"`
	DatabaseGrant       types.Set    `tfsdk:"database_grant"`
	ID                  types.String `tfsdk:"id"`
	PasswordCreated     types.String `tfsdk:"password_created"`
	TFAEnabled          types.Bool   `tfsdk:"tfa_enabled"`
	VerifiedPhoneNumber types.String `tfsdk:"verified_phone_number"`
}

type UserModel struct {
	Username            types.String       `tfsdk:"username"`
	SSHKeys             types.List         `tfsdk:"ssh_keys"`
	Email               types.String       `tfsdk:"email"`
	Restricted          types.Bool         `tfsdk:"restricted"`
	GlobalGrants        []GlobalGrantModel `tfsdk:"global_grants"`
	DomainGrant         []UserGrantModel   `tfsdk:"domain_grant"`
	FirewallGrant       []UserGrantModel   `tfsdk:"firewall_grant"`
	ImageGrant          []UserGrantModel   `tfsdk:"image_grant"`
	LinodeGrant         []UserGrantModel   `tfsdk:"linode_grant"`
	LongviewGrant       []UserGrantModel   `tfsdk:"longview_grant"`
	NodebalancerGrant   []UserGrantModel   `tfsdk:"nodebalancer_grant"`
	StackscriptGrant    []UserGrantModel   `tfsdk:"stackscript_grant"`
	VolumeGrant         []UserGrantModel   `tfsdk:"volume_grant"`
	DatabaseGrant       []UserGrantModel   `tfsdk:"database_grant"`
	ID                  types.String       `tfsdk:"id"`
	PasswordCreated     types.String       `tfsdk:"password_created"`
	TFAEnabled          types.Bool         `tfsdk:"tfa_enabled"`
	VerifiedPhoneNumber types.String       `tfsdk:"verified_phone_number"`
}

type GlobalGrantModel struct {
	AccountAccess        types.String `tfsdk:"account_access"`
	AddDatabases         types.Bool   `tfsdk:"add_databases"`
	AddDomains           types.Bool   `tfsdk:"add_domains"`
	AddFirewalls         types.Bool   `tfsdk:"add_firewalls"`
	AddImages            types.Bool   `tfsdk:"add_images"`
	AddLinodes           types.Bool   `tfsdk:"add_linodes"`
	AddLongview          types.Bool   `tfsdk:"add_longview"`
	AddNodebalancers     types.Bool   `tfsdk:"add_nodebalancers"`
	AddStackScripts      types.Bool   `tfsdk:"add_stackscripts"`
	AddVolumes           types.Bool   `tfsdk:"add_volumes"`
	CancelAccount        types.Bool   `tfsdk:"cancel_account"`
	LongviewSubscription types.Bool   `tfsdk:"longview_subscription"`
}

type UserGrantModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Permissions types.String `tfsdk:"permissions"`
	Label       types.String `tfsdk:"label"`
}

func (data *UserModel) parseComputedAttrs(
	ctx context.Context,
	user *linodego.User,
) diag.Diagnostics {
	sshKeys, diags := types.ListValueFrom(ctx, types.StringType, user.SSHKeys)
	if diags.HasError() {
		return diags
	}
	data.SSHKeys = sshKeys

	data.TFAEnabled = types.BoolValue(user.TFAEnabled)
	data.VerifiedPhoneNumber = types.StringPointerValue(user.VerifiedPhoneNumber)
	return nil
}

func (data *DataSourceModel) ParseUser(
	ctx context.Context, user *linodego.User,
) diag.Diagnostics {
	data.Username = types.StringValue(user.Username)
	data.Email = types.StringValue(user.Email)
	data.Restricted = types.BoolValue(user.Restricted)
	data.TFAEnabled = types.BoolValue(user.TFAEnabled)
	data.VerifiedPhoneNumber = types.StringPointerValue(user.VerifiedPhoneNumber)

	if user.PasswordCreated != nil {
		data.PasswordCreated = types.StringValue(user.PasswordCreated.Format(time.RFC3339))
	} else {
		data.PasswordCreated = types.StringNull()
	}

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

func (data *DataSourceModel) ParseUserGrants(
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

func (data *DataSourceModel) ParseNonUserGrants() {
	data.DatabaseGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.DomainGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.FirewallGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.GlobalGrants = types.ListNull(linodeUserGrantsGlobalObjectType)
	data.ImageGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.LinodeGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.LongviewGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.NodebalancerGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.StackscriptGrant = types.SetNull(linodeUserGrantsEntityObjectType)
	data.VolumeGrant = types.SetNull(linodeUserGrantsEntityObjectType)
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
	result["label"] = types.StringValue(string(entity.Label))

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

func expandGlobalGrant(global GlobalGrantModel) linodego.GlobalUserGrants {
	result := linodego.GlobalUserGrants{}

	result.AccountAccess = nil

	if !global.AccountAccess.IsNull() && !global.AccountAccess.IsUnknown() {
		accountAccess := linodego.GrantPermissionLevel(global.AccountAccess.ValueString())
		result.AccountAccess = &accountAccess
	}

	result.AddDomains = global.AddDomains.ValueBool()
	result.AddDatabases = global.AddDatabases.ValueBool()
	result.AddFirewalls = global.AddFirewalls.ValueBool()
	result.AddImages = global.AddImages.ValueBool()
	result.AddLinodes = global.AddLinodes.ValueBool()
	result.AddLongview = global.AddLongview.ValueBool()
	result.AddNodeBalancers = global.AddNodebalancers.ValueBool()
	result.AddStackScripts = global.AddStackScripts.ValueBool()
	result.AddVolumes = global.AddVolumes.ValueBool()
	result.CancelAccount = global.CancelAccount.ValueBool()
	result.LongviewSubscription = global.LongviewSubscription.ValueBool()

	return result
}

func expandUserGrantsEntities(entities []UserGrantModel) []linodego.EntityUserGrant {
	result := make([]linodego.EntityUserGrant, len(entities))

	for i, entity := range entities {
		userGrant := linodego.EntityUserGrant{}

		permissions := linodego.GrantPermissionLevel(entity.Permissions.ValueString())
		userGrant.ID = int(entity.ID.ValueInt64())
		userGrant.Permissions = &permissions
		result[i] = userGrant
	}

	return result
}
