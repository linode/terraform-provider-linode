package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

// GetFrameworkProviderModelFromSDKv2ProviderConfig is used by crossplane-provider-linode to convert from SDKv2 to FrameworkModel
func GetFrameworkProviderModelFromSDKv2ProviderConfig(config *Config) *FrameworkProviderModel {
	return &FrameworkProviderModel{
		AccessToken:                  types.StringValue(config.AccessToken),
		APIURL:                       types.StringValue(config.APIURL),
		APIVersion:                   types.StringValue(config.APIVersion),
		APICAPath:                    types.StringValue(config.APICAPath),
		UAPrefix:                     types.StringValue(config.UAPrefix),
		ConfigPath:                   types.StringValue(config.ConfigPath),
		ConfigProfile:                types.StringValue(config.ConfigProfile),
		SkipInstanceReadyPoll:        types.BoolValue(config.SkipInstanceReadyPoll),
		SkipInstanceDeletePoll:       types.BoolValue(config.SkipInstanceDeletePoll),
		SkipImplicitReboots:          types.BoolValue(config.SkipImplicitReboots),
		DisableInternalCache:         types.BoolValue(config.DisableInternalCache),
		MinRetryDelayMilliseconds:    types.Int64Value(int64(config.MinRetryDelayMilliseconds)),
		MaxRetryDelayMilliseconds:    types.Int64Value(int64(config.MaxRetryDelayMilliseconds)),
		EventPollMilliseconds:        types.Int64Value(int64(config.EventPollMilliseconds)),
		LKEEventPollMilliseconds:     types.Int64Value(int64(config.LKEEventPollMilliseconds)),
		LKENodeReadyPollMilliseconds: types.Int64Value(int64(config.LKENodeReadyPollMilliseconds)),
		ObjAccessKey:                 types.StringValue(config.ObjAccessKey),
		ObjSecretKey:                 types.StringValue(config.ObjSecretKey),
		ObjUseTempKeys:               types.BoolValue(config.ObjUseTempKeys),
		ObjBucketForceDelete:         types.BoolValue(config.ObjBucketForceDelete),
	}
}

type FrameworkProviderModel struct {
	AccessToken types.String `tfsdk:"token"`
	APIURL      types.String `tfsdk:"url"`
	APIVersion  types.String `tfsdk:"api_version"`
	APICAPath   types.String `tfsdk:"api_ca_path"`
	UAPrefix    types.String `tfsdk:"ua_prefix"`

	ConfigPath    types.String `tfsdk:"config_path"`
	ConfigProfile types.String `tfsdk:"config_profile"`

	SkipInstanceReadyPoll  types.Bool `tfsdk:"skip_instance_ready_poll"`
	SkipInstanceDeletePoll types.Bool `tfsdk:"skip_instance_delete_poll"`

	SkipImplicitReboots types.Bool `tfsdk:"skip_implicit_reboots"`

	DisableInternalCache types.Bool `tfsdk:"disable_internal_cache"`

	MinRetryDelayMilliseconds types.Int64 `tfsdk:"min_retry_delay_ms"`
	MaxRetryDelayMilliseconds types.Int64 `tfsdk:"max_retry_delay_ms"`

	EventPollMilliseconds    types.Int64 `tfsdk:"event_poll_ms"`
	LKEEventPollMilliseconds types.Int64 `tfsdk:"lke_event_poll_ms"`

	LKENodeReadyPollMilliseconds types.Int64 `tfsdk:"lke_node_ready_poll_ms"`

	ObjAccessKey         types.String `tfsdk:"obj_access_key"`
	ObjSecretKey         types.String `tfsdk:"obj_secret_key"`
	ObjUseTempKeys       types.Bool   `tfsdk:"obj_use_temp_keys"`
	ObjBucketForceDelete types.Bool   `tfsdk:"obj_bucket_force_delete"`
}

type FrameworkProviderMeta struct {
	Client *linodego.Client
	Config *FrameworkProviderModel
}
