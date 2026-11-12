//go:build unit

package nfsquota

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/require"
)

func ptrInt64(v int64) *int64 { return linodego.Pointer(v) }

func TestFlattenQuota_basic(t *testing.T) {
	quota := &linodego.NFSQuota{
		ID:                42,
		FilesystemID:      7,
		Path:              "/exports/data",
		IsRoot:            false,
		MaxCapacityBytes:  ptrInt64(2147483648),
		UsedCapacityBytes: ptrInt64(1024),
		QuotaExceeded:     false,
		MaxFileCount:      ptrInt64(2000),
		UsedFileCount:     ptrInt64(10),
		Status:            linodego.NFSQuotaStatusActive,
	}

	var model ResourceModel
	model.FlattenQuota(quota, false)

	require.Equal(t, "42", model.ID.ValueString())
	require.Equal(t, int64(7), model.FilesystemID.ValueInt64())
	require.Equal(t, "/exports/data", model.Path.ValueString())
	require.False(t, model.IsRoot.ValueBool())
	require.Equal(t, int64(2147483648), model.MaxCapacityBytes.ValueInt64())
	require.Equal(t, int64(1024), model.UsedCapacityBytes.ValueInt64())
	require.Equal(t, int64(2000), model.MaxFileCount.ValueInt64())
	require.Equal(t, int64(10), model.UsedFileCount.ValueInt64())
	require.Equal(t, "active", model.Status.ValueString())
	require.Nil(t, model.UserGroupConfig)
}

func TestFlattenQuota_userGroupConfig(t *testing.T) {
	quota := &linodego.NFSQuota{
		ID:           1,
		FilesystemID: 2,
		Path:         "/exports",
		Status:       linodego.NFSQuotaStatusActive,
		UserGroupConfig: &linodego.NFSUserGroupConfig{
			DefaultUserLimit: &linodego.NFSCapacityLimit{
				MaxCapacityBytes: ptrInt64(1073741824),
				MaxFileCount:     ptrInt64(1000),
			},
			UserLimits: []linodego.NFSIdentifiedLimit{
				{
					IdentifierType:   linodego.NFSQuotaRuleIdentifierTypeUID,
					Identifier:       "1001",
					MaxCapacityBytes: ptrInt64(536870912),
					MaxFileCount:     ptrInt64(500),
				},
			},
			UserGroupLimits: []linodego.NFSUserGroupCombinedLimit{
				{
					User: linodego.NFSUserGroupLimitIdentity{
						IdentifierType: linodego.NFSQuotaRuleIdentifierTypeUID,
						Identifier:     "1001",
					},
					Group: linodego.NFSUserGroupLimitIdentity{
						IdentifierType: linodego.NFSQuotaRuleIdentifierTypeGID,
						Identifier:     "2001",
					},
					MaxCapacityBytes: ptrInt64(268435456),
					MaxFileCount:     ptrInt64(250),
				},
			},
		},
	}

	var model ResourceModel
	model.FlattenQuota(quota, false)

	require.NotNil(t, model.UserGroupConfig)
	require.NotNil(t, model.UserGroupConfig.DefaultUserLimit)
	require.Equal(t, int64(1073741824), model.UserGroupConfig.DefaultUserLimit.MaxCapacityBytes.ValueInt64())
	require.Nil(t, model.UserGroupConfig.DefaultGroupLimit)

	require.Len(t, model.UserGroupConfig.UserLimits, 1)
	require.Equal(t, "uid", model.UserGroupConfig.UserLimits[0].IdentifierType.ValueString())
	require.Equal(t, "1001", model.UserGroupConfig.UserLimits[0].Identifier.ValueString())

	require.Len(t, model.UserGroupConfig.UserGroupLimits, 1)
	require.Equal(t, "gid", model.UserGroupConfig.UserGroupLimits[0].Group.IdentifierType.ValueString())
	require.Equal(t, int64(268435456), model.UserGroupConfig.UserGroupLimits[0].MaxCapacityBytes.ValueInt64())
}

func TestFlattenQuota_nullPointers(t *testing.T) {
	quota := &linodego.NFSQuota{
		ID:           5,
		FilesystemID: 6,
		Path:         "/exports",
		Status:       linodego.NFSQuotaStatusUpdating,
	}

	var model ResourceModel
	model.FlattenQuota(quota, false)

	require.True(t, model.MaxCapacityBytes.IsNull())
	require.True(t, model.UsedCapacityBytes.IsNull())
	require.True(t, model.MaxFileCount.IsNull())
	require.True(t, model.UsedFileCount.IsNull())
}

func TestToCreateOptions(t *testing.T) {
	model := ResourceModel{
		Path:             types.StringValue("/exports/create"),
		MaxCapacityBytes: types.Int64Value(2147483648),
		MaxFileCount:     types.Int64Value(2000),
		UserGroupConfig: &UserGroupConfigModel{
			DefaultUserLimit: &CapacityLimitModel{
				MaxCapacityBytes: types.Int64Value(1073741824),
				MaxFileCount:     types.Int64Value(1000),
			},
			UserLimits: []IdentifiedLimitModel{
				{
					IdentifierType:   types.StringValue("uid"),
					Identifier:       types.StringValue("1001"),
					MaxCapacityBytes: types.Int64Value(536870912),
					MaxFileCount:     types.Int64Value(500),
				},
			},
		},
	}

	opts := model.ToCreateOptions()

	require.Equal(t, "/exports/create", opts.Path)
	require.Equal(t, int64(2147483648), opts.MaxCapacityBytes)
	require.Equal(t, int64(2000), opts.MaxFileCount)
	require.NotNil(t, opts.UserGroupConfig)
	require.NotNil(t, opts.UserGroupConfig.DefaultUserLimit)
	require.Equal(t, int64(1073741824), *opts.UserGroupConfig.DefaultUserLimit.MaxCapacityBytes)
	require.NotNil(t, opts.UserGroupConfig.UserLimits)
	require.Len(t, *opts.UserGroupConfig.UserLimits, 1)
	require.Equal(t, linodego.NFSQuotaRuleIdentifierTypeUID, (*opts.UserGroupConfig.UserLimits)[0].IdentifierType)
}

func TestParseQuota(t *testing.T) {
	quota := &linodego.NFSQuota{
		ID:               9,
		FilesystemID:     3,
		Path:             "/exports/ds",
		MaxCapacityBytes: ptrInt64(2147483648),
		MaxFileCount:     ptrInt64(2000),
		Status:           linodego.NFSQuotaStatusActive,
	}

	var model DataSourceModel
	model.ParseQuota(quota)

	require.Equal(t, int64(9), model.ID.ValueInt64())
	require.Equal(t, int64(3), model.FilesystemID.ValueInt64())
	require.Equal(t, "/exports/ds", model.Path.ValueString())
	require.Equal(t, int64(2147483648), model.MaxCapacityBytes.ValueInt64())
	require.Equal(t, "active", model.Status.ValueString())
}
