package linode

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func (fp *FrameworkProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data helper.FrameworkProviderModel
	var meta helper.FrameworkProviderMeta

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fp.HandleDefaults(&data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	fp.InitProvider(ctx, &data, req.TerraformVersion, &resp.Diagnostics, &meta)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = &meta
	resp.DataSourceData = &meta

	fp.Meta = &meta
}

// We should replace this with an official validator if
// HashiCorp decide to implement it in the future
// feature request track:
// https://github.com/hashicorp/terraform-plugin-framework-validators/issues/125
func (fp *FrameworkProvider) ValidateConfig(
	ctx context.Context,
	req provider.ValidateConfigRequest,
	resp *provider.ValidateConfigResponse,
) {
	var data helper.FrameworkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := url.Parse(data.APIURL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse the base API URL in the configuration",
			err.Error(),
		)
	}
}

func GetIntFromEnv(
	key string,
	defaultValue basetypes.Int64Value,
	diags *diag.Diagnostics,
) basetypes.Int64Value {
	envVarVal := os.Getenv(key)
	if envVarVal == "" {
		return defaultValue
	}

	intVal, err := strconv.ParseInt(envVarVal, 10, 64)
	if err != nil {
		diags.AddWarning(
			fmt.Sprintf(
				"Failed to parse the environment variable %v "+
					"to an integer. Will use default value: %v instead",
				key,
				defaultValue.ValueInt64(),
			),
			err.Error(),
		)

		return defaultValue
	}

	return types.Int64Value(intVal)
}

func GetStringFromEnv(key string, defaultValue basetypes.StringValue) basetypes.StringValue {
	envVarVal := os.Getenv(key)

	if envVarVal == "" {
		return defaultValue
	}

	return types.StringValue(envVarVal)
}

func (fp *FrameworkProvider) HandleDefaults(
	lpm *helper.FrameworkProviderModel,
	diags *diag.Diagnostics,
) {
	if lpm.AccessToken.IsNull() {
		lpm.AccessToken = GetStringFromEnv("LINODE_TOKEN", types.StringNull())
	}

	if lpm.ConfigPath.IsNull() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			diags.AddError(
				"Failed to get the home directory of the user for the config path.",
				err.Error(),
			)
			return
		}
		configPath := fmt.Sprintf("%s/.config/linode", homeDir)
		lpm.ConfigPath = types.StringValue(configPath)
	}

	if lpm.ConfigProfile.IsNull() {
		lpm.ConfigProfile = types.StringValue("default")
	}

	if lpm.APIURL.IsNull() {
		lpm.APIURL = GetStringFromEnv(
			"LINODE_URL",
			types.StringValue(helper.DefaultLinodeURL),
		)
	}

	if lpm.UAPrefix.IsNull() {
		lpm.UAPrefix = GetStringFromEnv("LINODE_UA_PREFIX", types.StringNull())
	}

	if lpm.APIVersion.IsNull() {
		lpm.APIVersion = GetStringFromEnv(
			"LINODE_API_VERSION",
			types.StringValue("v4"),
		)
	}

	if lpm.SkipInstanceReadyPoll.IsNull() {
		lpm.SkipInstanceReadyPoll = types.BoolValue(false)
	}

	if lpm.SkipInstanceDeletePoll.IsNull() {
		lpm.SkipInstanceDeletePoll = types.BoolValue(false)
	}

	if lpm.SkipImplicitReboots.IsNull() {
		lpm.SkipImplicitReboots = types.BoolValue(false)
	}

	if lpm.DisableInternalCache.IsNull() {
		lpm.DisableInternalCache = types.BoolValue(false)
	}

	if lpm.EventPollMilliseconds.IsNull() {
		lpm.EventPollMilliseconds = GetIntFromEnv(
			"LINODE_EVENT_POLL_MS",
			types.Int64Value(4000),
			diags,
		)
	}

	if lpm.LKEEventPollMilliseconds.IsNull() {
		lpm.LKEEventPollMilliseconds = types.Int64Value(3000)
	}

	if lpm.LKENodeReadyPollMilliseconds.IsNull() {
		lpm.LKENodeReadyPollMilliseconds = types.Int64Value(3000)
	}

	if lpm.ObjAccessKey.IsNull() {
		lpm.ObjAccessKey = GetStringFromEnv(
			"LINODE_OBJ_ACCESS_KEY",
			types.StringNull(),
		)
	}

	if lpm.ObjSecretKey.IsNull() {
		lpm.ObjSecretKey = GetStringFromEnv(
			"LINODE_OBJ_SECRET_KEY",
			types.StringNull(),
		)
	}

	if lpm.ObjUseTempKeys.IsNull() {
		lpm.ObjUseTempKeys = types.BoolValue(false)
	}
}

func (fp *FrameworkProvider) InitProvider(
	ctx context.Context,
	lpm *helper.FrameworkProviderModel,
	tfVersion string,
	diags *diag.Diagnostics,
	meta *helper.FrameworkProviderMeta,
) {
	if fp.Meta.Client != nil {
		// Crossplane provider-linode expects to use a single configured instance of the linode client across all invocations
		// However, due to how upjet operates, the configureProvider() gets invoked on every resource call. To preserve the client,
		// see if the fp.Meta.Client is already initialized, and if so, re-use it.

		tflog.Info(ctx, "Linode client was already configured, re-using..")
		meta.Client = fp.Meta.Client
		meta.Config = fp.Meta.Config
		return
	}

	loggingTransport := helper.NewAPILoggerTransport(
		logging.NewSubsystemLoggingHTTPTransport(
			helper.APILoggerSubsystem,
			http.DefaultTransport,
		),
	)

	oauth2Client := &http.Client{
		Transport: loggingTransport,
	}

	accessToken := lpm.AccessToken.ValueString()
	APIURL := lpm.APIURL.ValueString()
	APIVersion := lpm.APIVersion.ValueString()
	UAPrefix := lpm.UAPrefix.ValueString()

	configPath := lpm.ConfigPath.ValueString()
	configProfile := lpm.ConfigProfile.ValueString()

	// skipInstanceReadyPoll := lpm.SkipInstanceReadyPoll.ValueBool()
	// skipInstanceDeletePoll := lpm.SkipInstanceDeletePoll.ValueBool()

	disableInternalCache := lpm.DisableInternalCache.ValueBool()

	minRetryDelayMilliseconds := lpm.MinRetryDelayMilliseconds.ValueInt64()
	maxRetryDelayMilliseconds := lpm.MaxRetryDelayMilliseconds.ValueInt64()

	eventPollMilliseconds := lpm.EventPollMilliseconds.ValueInt64()
	// LKENodeReadyPollMilliseconds := lpm.LKEEventPollMilliseconds.ValueInt64()

	client := linodego.NewClient(oauth2Client)
	// Load the config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		tflog.Info(ctx, "Using Linode profile", map[string]any{
			"config_path": lpm.ConfigPath,
		})
		err = client.LoadConfig(&linodego.LoadConfigOptions{
			Path:    configPath,
			Profile: configProfile,
		})
		if err != nil {
			diags.AddError("Error occurs when loading linode profile.", err.Error())
			return
		}
	} else {
		tflog.Info(ctx, "Linode config does not exist, skipping..")
	}

	// Overrides
	if accessToken != "" {
		client.SetToken(accessToken)
	}

	if APIURL != "" {
		client.SetBaseURL(APIURL)
	}

	if len(APIVersion) > 0 {
		client.SetAPIVersion(APIVersion)
	}

	client.UseCache(!disableInternalCache)

	if eventPollMilliseconds != 0 {
		client.SetPollDelay(time.Duration(eventPollMilliseconds) * time.Millisecond)
	}

	if minRetryDelayMilliseconds != 0 {
		client.SetRetryWaitTime(time.Duration(minRetryDelayMilliseconds) * time.Millisecond)
	}
	if maxRetryDelayMilliseconds != 0 {
		client.SetRetryMaxWaitTime(time.Duration(maxRetryDelayMilliseconds) * time.Millisecond)
	}

	userAgent := fp.terraformUserAgent(tfVersion, UAPrefix)
	client.SetUserAgent(userAgent)

	helper.ApplyAllRetryConditions(&client)

	meta.Config = lpm
	meta.Client = &client
}

func (fp *FrameworkProvider) terraformUserAgent(
	tfVersion string,
	UAPrefix string,
) string {
	userAgent := strings.TrimSpace(
		fmt.Sprintf(
			"HashiCorp Terraform/%s (+https://www.terraform.io) "+
				"Terraform-Plugin-Framework/%s terraform-provider-linode/%s",
			tfVersion, helper.GetFrameworkVersion(), fp.ProviderVersion,
		),
	)

	if add := os.Getenv(helper.UAEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			userAgent += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", userAgent)
		}
	}

	if UAPrefix != "" {
		userAgent = UAPrefix + " " + userAgent
	}

	return userAgent
}
