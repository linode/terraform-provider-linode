package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/terraform-provider-linode/linode/account"
	"github.com/linode/terraform-provider-linode/linode/backup"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/instancenetworking"
	"github.com/linode/terraform-provider-linode/linode/kernel"
	"github.com/linode/terraform-provider-linode/linode/lkeversions"
	"github.com/linode/terraform-provider-linode/linode/networkingip"
	"github.com/linode/terraform-provider-linode/linode/profile"
	"github.com/linode/terraform-provider-linode/linode/stackscript"
	"github.com/linode/terraform-provider-linode/linode/token"
)

type FrameworkProvider struct {
	ProviderVersion string
	Meta            *helper.FrameworkProviderMeta
}

func CreateFrameworkProvider(version string) provider.ProviderWithValidateConfig {
	return &FrameworkProvider{
		ProviderVersion: version,
	}
}

func (p *FrameworkProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "linodecloud"
}

func (p *FrameworkProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "The token that allows you access to your Linode account",
			},
			"config_path": schema.StringAttribute{
				Optional: true,
			},
			"config_profile": schema.StringAttribute{
				Optional: true,
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The HTTP(S) API address of the Linode API to use.",
			},
			"ua_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "An HTTP User-Agent Prefix to prepend in API requests.",
			},
			"api_version": schema.StringAttribute{
				Optional:    true,
				Description: "The version of Linode API.",
			},
			"skip_instance_ready_poll": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip waiting for a linode_instance resource to be running.",
			},
			"skip_instance_delete_poll": schema.BoolAttribute{
				Optional:    true,
				Description: "Skip waiting for a linode_instance resource to finish deleting.",
			},
			"disable_internal_cache": schema.BoolAttribute{
				Optional:    true,
				Description: "Disable the internal caching system that backs certain Linode API requests.",
			},
			"min_retry_delay_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum delay in milliseconds before retrying a request.",
			},
			"max_retry_delay_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum delay in milliseconds before retrying a request.",
			},
			"event_poll_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "The rate in milliseconds to poll for events.",
			},
			"lke_event_poll_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "The rate in milliseconds to poll for LKE events.",
			},
			"lke_node_ready_poll_ms": schema.Int64Attribute{
				Optional:    true,
				Description: "The rate in milliseconds to poll for an LKE node to be ready.",
			},
		},
	}
}

func (p *FrameworkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		token.NewResource,
		stackscript.NewResource,
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		account.NewDataSource,
		backup.NewDataSource,
		kernel.NewDataSource,
		stackscript.NewDataSource,
		profile.NewDataSource,
		networkingip.NewDataSource,
		lkeversions.NewDataSource,
		instancenetworking.NewDataSource,
	}
}
