package nfsquota

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

// ResourceModel is the Terraform state model for an NFS Quota resource.
type ResourceModel struct {
	ID                types.String          `tfsdk:"id"`
	SpaceID           types.Int64           `tfsdk:"space_id"`
	FilesystemID      types.Int64           `tfsdk:"filesystem_id"`
	Path              types.String          `tfsdk:"path"`
	IsRoot            types.Bool            `tfsdk:"is_root"`
	MaxCapacityBytes  types.Int64           `tfsdk:"max_capacity_bytes"`
	UsedCapacityBytes types.Int64           `tfsdk:"used_capacity_bytes"`
	QuotaExceeded     types.Bool            `tfsdk:"quota_exceeded"`
	MaxFileCount      types.Int64           `tfsdk:"max_file_count"`
	UsedFileCount     types.Int64           `tfsdk:"used_file_count"`
	Status            types.String          `tfsdk:"status"`
	Created           timetypes.RFC3339     `tfsdk:"created"`
	Updated           timetypes.RFC3339     `tfsdk:"updated"`
	CollectedAt       timetypes.RFC3339     `tfsdk:"collected_at"`
	UserGroupConfig   *UserGroupConfigModel `tfsdk:"user_group_config"`
}

// UserGroupConfigModel models user and group quota configuration.
type UserGroupConfigModel struct {
	DefaultUserLimit  *CapacityLimitModel           `tfsdk:"default_user_limit"`
	DefaultGroupLimit *CapacityLimitModel           `tfsdk:"default_group_limit"`
	UserLimits        []IdentifiedLimitModel        `tfsdk:"user_limits"`
	GroupLimits       []IdentifiedLimitModel        `tfsdk:"group_limits"`
	UserGroupLimits   []UserGroupCombinedLimitModel `tfsdk:"user_group_limits"`
}

// CapacityLimitModel models capacity and file count limits.
type CapacityLimitModel struct {
	MaxCapacityBytes types.Int64 `tfsdk:"max_capacity_bytes"`
	MaxFileCount     types.Int64 `tfsdk:"max_file_count"`
}

// IdentifiedLimitModel models limits for a specific user or group.
type IdentifiedLimitModel struct {
	IdentifierType   types.String `tfsdk:"identifier_type"`
	Identifier       types.String `tfsdk:"identifier"`
	MaxCapacityBytes types.Int64  `tfsdk:"max_capacity_bytes"`
	MaxFileCount     types.Int64  `tfsdk:"max_file_count"`
}

// UserGroupCombinedLimitModel models limits for a combined user and group.
type UserGroupCombinedLimitModel struct {
	User             IdentityModel `tfsdk:"user"`
	Group            IdentityModel `tfsdk:"group"`
	MaxCapacityBytes types.Int64   `tfsdk:"max_capacity_bytes"`
	MaxFileCount     types.Int64   `tfsdk:"max_file_count"`
}

// IdentityModel identifies a user or group.
type IdentityModel struct {
	IdentifierType types.String `tfsdk:"identifier_type"`
	Identifier     types.String `tfsdk:"identifier"`
}

// FlattenQuota populates the resource model from an API response.
func (m *ResourceModel) FlattenQuota(quota *linodego.NFSQuota, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateString(m.ID, strconv.Itoa(quota.ID), preserveKnown)
	m.FilesystemID = helper.KeepOrUpdateInt64(m.FilesystemID, int64(quota.FilesystemID), preserveKnown)
	m.Path = helper.KeepOrUpdateString(m.Path, quota.Path, preserveKnown)
	m.IsRoot = helper.KeepOrUpdateBool(m.IsRoot, quota.IsRoot, preserveKnown)
	m.MaxCapacityBytes = helper.KeepOrUpdateInt64Pointer(m.MaxCapacityBytes, quota.MaxCapacityBytes, preserveKnown)
	m.UsedCapacityBytes = helper.KeepOrUpdateInt64Pointer(m.UsedCapacityBytes, quota.UsedCapacityBytes, preserveKnown)
	m.QuotaExceeded = helper.KeepOrUpdateBool(m.QuotaExceeded, quota.QuotaExceeded, preserveKnown)
	m.MaxFileCount = helper.KeepOrUpdateInt64Pointer(m.MaxFileCount, quota.MaxFileCount, preserveKnown)
	m.UsedFileCount = helper.KeepOrUpdateInt64Pointer(m.UsedFileCount, quota.UsedFileCount, preserveKnown)
	m.Status = helper.KeepOrUpdateString(m.Status, string(quota.Status), preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, timetypes.NewRFC3339TimePointerValue(quota.Created), preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, timetypes.NewRFC3339TimePointerValue(quota.Updated), preserveKnown)
	m.CollectedAt = helper.KeepOrUpdateValue(m.CollectedAt, timetypes.NewRFC3339TimePointerValue(quota.CollectedAt), preserveKnown)

	// user_group_config is user-configured. During create/update (preserveKnown)
	// the plan is authoritative, so we only refresh it from the API on reads.
	if !preserveKnown {
		m.UserGroupConfig = flattenUserGroupConfig(quota.UserGroupConfig)
	}
}

// CopyFrom copies values from another model instance, respecting preserveKnown.
func (m *ResourceModel) CopyFrom(other ResourceModel, preserveKnown bool) {
	m.ID = helper.KeepOrUpdateValue(m.ID, other.ID, preserveKnown)
	m.SpaceID = helper.KeepOrUpdateValue(m.SpaceID, other.SpaceID, preserveKnown)
	m.FilesystemID = helper.KeepOrUpdateValue(m.FilesystemID, other.FilesystemID, preserveKnown)
	m.Path = helper.KeepOrUpdateValue(m.Path, other.Path, preserveKnown)
	m.IsRoot = helper.KeepOrUpdateValue(m.IsRoot, other.IsRoot, preserveKnown)
	m.MaxCapacityBytes = helper.KeepOrUpdateValue(m.MaxCapacityBytes, other.MaxCapacityBytes, preserveKnown)
	m.UsedCapacityBytes = helper.KeepOrUpdateValue(m.UsedCapacityBytes, other.UsedCapacityBytes, preserveKnown)
	m.QuotaExceeded = helper.KeepOrUpdateValue(m.QuotaExceeded, other.QuotaExceeded, preserveKnown)
	m.MaxFileCount = helper.KeepOrUpdateValue(m.MaxFileCount, other.MaxFileCount, preserveKnown)
	m.UsedFileCount = helper.KeepOrUpdateValue(m.UsedFileCount, other.UsedFileCount, preserveKnown)
	m.Status = helper.KeepOrUpdateValue(m.Status, other.Status, preserveKnown)
	m.Created = helper.KeepOrUpdateValue(m.Created, other.Created, preserveKnown)
	m.Updated = helper.KeepOrUpdateValue(m.Updated, other.Updated, preserveKnown)
	m.CollectedAt = helper.KeepOrUpdateValue(m.CollectedAt, other.CollectedAt, preserveKnown)

	if m.UserGroupConfig == nil && other.UserGroupConfig != nil {
		copied := *other.UserGroupConfig
		m.UserGroupConfig = &copied
	}
}

// ToCreateOptions builds the linodego create options from the resource model.
func (m *ResourceModel) ToCreateOptions() linodego.NFSQuotaCreateOptions {
	opts := linodego.NFSQuotaCreateOptions{
		Path:             m.Path.ValueString(),
		MaxCapacityBytes: m.MaxCapacityBytes.ValueInt64(),
		MaxFileCount:     m.MaxFileCount.ValueInt64(),
	}

	if m.UserGroupConfig != nil {
		opts.UserGroupConfig = m.UserGroupConfig.toOptions()
	}

	return opts
}

// DataSourceModel is the Terraform state model for an NFS Quota data source.
type DataSourceModel struct {
	ID                types.Int64           `tfsdk:"id"`
	SpaceID           types.Int64           `tfsdk:"space_id"`
	FilesystemID      types.Int64           `tfsdk:"filesystem_id"`
	Path              types.String          `tfsdk:"path"`
	IsRoot            types.Bool            `tfsdk:"is_root"`
	MaxCapacityBytes  types.Int64           `tfsdk:"max_capacity_bytes"`
	UsedCapacityBytes types.Int64           `tfsdk:"used_capacity_bytes"`
	QuotaExceeded     types.Bool            `tfsdk:"quota_exceeded"`
	MaxFileCount      types.Int64           `tfsdk:"max_file_count"`
	UsedFileCount     types.Int64           `tfsdk:"used_file_count"`
	Status            types.String          `tfsdk:"status"`
	Created           timetypes.RFC3339     `tfsdk:"created"`
	Updated           timetypes.RFC3339     `tfsdk:"updated"`
	CollectedAt       timetypes.RFC3339     `tfsdk:"collected_at"`
	UserGroupConfig   *UserGroupConfigModel `tfsdk:"user_group_config"`
}

// ParseQuota populates the data source model from an API response.
func (m *DataSourceModel) ParseQuota(quota *linodego.NFSQuota) {
	m.ID = types.Int64Value(int64(quota.ID))
	m.FilesystemID = types.Int64Value(int64(quota.FilesystemID))
	m.Path = types.StringValue(quota.Path)
	m.IsRoot = types.BoolValue(quota.IsRoot)
	m.MaxCapacityBytes = int64PointerValue(quota.MaxCapacityBytes)
	m.UsedCapacityBytes = int64PointerValue(quota.UsedCapacityBytes)
	m.QuotaExceeded = types.BoolValue(quota.QuotaExceeded)
	m.MaxFileCount = int64PointerValue(quota.MaxFileCount)
	m.UsedFileCount = int64PointerValue(quota.UsedFileCount)
	m.Status = types.StringValue(string(quota.Status))
	m.Created = timetypes.NewRFC3339TimePointerValue(quota.Created)
	m.Updated = timetypes.NewRFC3339TimePointerValue(quota.Updated)
	m.CollectedAt = timetypes.NewRFC3339TimePointerValue(quota.CollectedAt)
	m.UserGroupConfig = flattenUserGroupConfig(quota.UserGroupConfig)
}

func int64PointerValue(v *int64) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(*v)
}

func flattenUserGroupConfig(config *linodego.NFSUserGroupConfig) *UserGroupConfigModel {
	if config == nil {
		return nil
	}

	result := &UserGroupConfigModel{
		DefaultUserLimit:  flattenCapacityLimit(config.DefaultUserLimit),
		DefaultGroupLimit: flattenCapacityLimit(config.DefaultGroupLimit),
	}

	if config.UserLimits != nil {
		result.UserLimits = helper.MapSlice(config.UserLimits, flattenIdentifiedLimit)
	}
	if config.GroupLimits != nil {
		result.GroupLimits = helper.MapSlice(config.GroupLimits, flattenIdentifiedLimit)
	}
	if config.UserGroupLimits != nil {
		result.UserGroupLimits = helper.MapSlice(config.UserGroupLimits, flattenUserGroupCombinedLimit)
	}

	return result
}

func flattenCapacityLimit(limit *linodego.NFSCapacityLimit) *CapacityLimitModel {
	if limit == nil {
		return nil
	}
	return &CapacityLimitModel{
		MaxCapacityBytes: int64PointerValue(limit.MaxCapacityBytes),
		MaxFileCount:     int64PointerValue(limit.MaxFileCount),
	}
}

func flattenIdentifiedLimit(limit linodego.NFSIdentifiedLimit) IdentifiedLimitModel {
	return IdentifiedLimitModel{
		IdentifierType:   types.StringValue(string(limit.IdentifierType)),
		Identifier:       types.StringValue(limit.Identifier),
		MaxCapacityBytes: int64PointerValue(limit.MaxCapacityBytes),
		MaxFileCount:     int64PointerValue(limit.MaxFileCount),
	}
}

func flattenUserGroupCombinedLimit(limit linodego.NFSUserGroupCombinedLimit) UserGroupCombinedLimitModel {
	return UserGroupCombinedLimitModel{
		User:             flattenIdentity(limit.User),
		Group:            flattenIdentity(limit.Group),
		MaxCapacityBytes: int64PointerValue(limit.MaxCapacityBytes),
		MaxFileCount:     int64PointerValue(limit.MaxFileCount),
	}
}

func flattenIdentity(identity linodego.NFSUserGroupLimitIdentity) IdentityModel {
	return IdentityModel{
		IdentifierType: types.StringValue(string(identity.IdentifierType)),
		Identifier:     types.StringValue(identity.Identifier),
	}
}

func (m *UserGroupConfigModel) toOptions() *linodego.NFSUserGroupConfigOptions {
	opts := &linodego.NFSUserGroupConfigOptions{
		DefaultUserLimit:  m.DefaultUserLimit.toOptions(),
		DefaultGroupLimit: m.DefaultGroupLimit.toOptions(),
	}

	if m.UserLimits != nil {
		limits := helper.MapSlice(m.UserLimits, identifiedLimitToOptions)
		opts.UserLimits = &limits
	}
	if m.GroupLimits != nil {
		limits := helper.MapSlice(m.GroupLimits, identifiedLimitToOptions)
		opts.GroupLimits = &limits
	}
	if m.UserGroupLimits != nil {
		limits := helper.MapSlice(m.UserGroupLimits, userGroupCombinedLimitToOptions)
		opts.UserGroupLimits = &limits
	}

	return opts
}

func (m *CapacityLimitModel) toOptions() *linodego.NFSCapacityLimitOptions {
	if m == nil {
		return nil
	}
	return &linodego.NFSCapacityLimitOptions{
		MaxCapacityBytes: m.MaxCapacityBytes.ValueInt64Pointer(),
		MaxFileCount:     m.MaxFileCount.ValueInt64Pointer(),
	}
}

func identifiedLimitToOptions(m IdentifiedLimitModel) linodego.NFSIdentifiedLimitOptions {
	return linodego.NFSIdentifiedLimitOptions{
		IdentifierType:   linodego.NFSQuotaRuleIdentifierType(m.IdentifierType.ValueString()),
		Identifier:       m.Identifier.ValueString(),
		MaxCapacityBytes: m.MaxCapacityBytes.ValueInt64Pointer(),
		MaxFileCount:     m.MaxFileCount.ValueInt64Pointer(),
	}
}

func userGroupCombinedLimitToOptions(m UserGroupCombinedLimitModel) linodego.NFSUserGroupCombinedLimitOptions {
	return linodego.NFSUserGroupCombinedLimitOptions{
		User:             identityToOptions(m.User),
		Group:            identityToOptions(m.Group),
		MaxCapacityBytes: m.MaxCapacityBytes.ValueInt64Pointer(),
		MaxFileCount:     m.MaxFileCount.ValueInt64Pointer(),
	}
}

func identityToOptions(m IdentityModel) linodego.NFSUserGroupLimitIdentity {
	return linodego.NFSUserGroupLimitIdentity{
		IdentifierType: linodego.NFSQuotaRuleIdentifierType(m.IdentifierType.ValueString()),
		Identifier:     m.Identifier.ValueString(),
	}
}
