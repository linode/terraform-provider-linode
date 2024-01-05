package linode

import (
	"context"

	"github.com/linode/terraform-provider-linode/v2/linode/accountavailabilities"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/account"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailability"
	"github.com/linode/terraform-provider-linode/v2/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/v2/linode/accountlogins"
	"github.com/linode/terraform-provider-linode/v2/linode/accountsettings"
	"github.com/linode/terraform-provider-linode/v2/linode/backup"
	"github.com/linode/terraform-provider-linode/v2/linode/databasebackups"
	"github.com/linode/terraform-provider-linode/v2/linode/databaseengines"
	"github.com/linode/terraform-provider-linode/v2/linode/databasemysql"
	"github.com/linode/terraform-provider-linode/v2/linode/databasepostgresql"
	"github.com/linode/terraform-provider-linode/v2/linode/databases"
	"github.com/linode/terraform-provider-linode/v2/linode/domain"
	"github.com/linode/terraform-provider-linode/v2/linode/domainrecord"
	"github.com/linode/terraform-provider-linode/v2/linode/domains"
	"github.com/linode/terraform-provider-linode/v2/linode/domainzonefile"
	"github.com/linode/terraform-provider-linode/v2/linode/firewall"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalldevice"
	"github.com/linode/terraform-provider-linode/v2/linode/firewalls"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/image"
	"github.com/linode/terraform-provider-linode/v2/linode/images"
	"github.com/linode/terraform-provider-linode/v2/linode/instanceip"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetype"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetypes"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6range"
	"github.com/linode/terraform-provider-linode/v2/linode/kernel"
	"github.com/linode/terraform-provider-linode/v2/linode/kernels"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversions"
	"github.com/linode/terraform-provider-linode/v2/linode/nb"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfigs"
	"github.com/linode/terraform-provider-linode/v2/linode/nbnode"
	"github.com/linode/terraform-provider-linode/v2/linode/nbs"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket"
	"github.com/linode/terraform-provider-linode/v2/linode/objcluster"
	"github.com/linode/terraform-provider-linode/v2/linode/objkey"
	"github.com/linode/terraform-provider-linode/v2/linode/profile"
	"github.com/linode/terraform-provider-linode/v2/linode/rdns"
	"github.com/linode/terraform-provider-linode/v2/linode/region"
	"github.com/linode/terraform-provider-linode/v2/linode/regions"
	"github.com/linode/terraform-provider-linode/v2/linode/sshkey"
	"github.com/linode/terraform-provider-linode/v2/linode/sshkeys"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscript"
	"github.com/linode/terraform-provider-linode/v2/linode/stackscripts"
	"github.com/linode/terraform-provider-linode/v2/linode/token"
	"github.com/linode/terraform-provider-linode/v2/linode/user"
	"github.com/linode/terraform-provider-linode/v2/linode/users"
	"github.com/linode/terraform-provider-linode/v2/linode/vlan"
	"github.com/linode/terraform-provider-linode/v2/linode/volume"
	"github.com/linode/terraform-provider-linode/v2/linode/volumes"
	"github.com/linode/terraform-provider-linode/v2/linode/vpc"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcs"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnet"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnets"
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
			"skip_implicit_reboots": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, Linode Instances will not be rebooted on config and interface changes.",
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
		rdns.NewResource,
		objkey.NewResource,
		sshkey.NewResource,
		ipv6range.NewResource,
		nb.NewResource,
		accountsettings.NewResource,
		vpcsubnet.NewResource,
		vpc.NewResource,
		instanceip.NewResource,
		firewalldevice.NewResource,
	}
}

func (p *FrameworkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		accountavailabilities.NewDataSource,
		account.NewDataSource,
		backup.NewDataSource,
		firewall.NewDataSource,
		kernel.NewDataSource,
		stackscript.NewDataSource,
		stackscripts.NewDataSource,
		profile.NewDataSource,
		nb.NewDataSource,
		networkingip.NewDataSource,
		lkeversions.NewDataSource,
		regions.NewDataSource,
		ipv6range.NewDataSource,
		objbucket.NewDataSource,
		sshkey.NewDataSource,
		sshkeys.NewDataSource,
		instancenetworking.NewDataSource,
		objcluster.NewDataSource,
		domainrecord.NewDataSource,
		databasepostgresql.NewDataSource,
		volume.NewDataSource,
		databasemysql.NewDataSource,
		domainzonefile.NewDataSource,
		domain.NewDataSource,
		user.NewDataSource,
		nbconfig.NewDataSource,
		instancetype.NewDataSource,
		instancetypes.NewDataSource,
		image.NewDataSource,
		images.NewDataSource,
		accountlogin.NewDataSource,
		accountlogins.NewDataSource,
		databasebackups.NewDataSource,
		databases.NewDataSource,
		databaseengines.NewDataSource,
		region.NewDataSource,
		vlan.NewDataSource,
		users.NewDataSource,
		nbnode.NewDataSource,
		nbs.NewDataSource,
		accountsettings.NewDataSource,
		firewalls.NewDataSource,
		kernels.NewDataSource,
		vpcsubnet.NewDataSource,
		vpc.NewDataSource,
		vpcsubnets.NewDataSource,
		vpcs.NewDataSource,
		volumes.NewDataSource,
		accountavailability.NewDataSource,
		nbconfigs.NewDataSource,
		domains.NewDataSource,
	}
}
