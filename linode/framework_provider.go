package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/account"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailabilities"
	"github.com/linode/terraform-provider-linode/v2/linode/accountavailability"
	"github.com/linode/terraform-provider-linode/v2/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/v2/linode/accountlogins"
	"github.com/linode/terraform-provider-linode/v2/linode/accountsettings"
	"github.com/linode/terraform-provider-linode/v2/linode/backup"
	"github.com/linode/terraform-provider-linode/v2/linode/childaccount"
	"github.com/linode/terraform-provider-linode/v2/linode/childaccounts"
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
	"github.com/linode/terraform-provider-linode/v2/linode/instancedisk"
	"github.com/linode/terraform-provider-linode/v2/linode/instanceip"
	"github.com/linode/terraform-provider-linode/v2/linode/instancenetworking"
	"github.com/linode/terraform-provider-linode/v2/linode/instancereservedip"
	"github.com/linode/terraform-provider-linode/v2/linode/instancesharedips"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetype"
	"github.com/linode/terraform-provider-linode/v2/linode/instancetypes"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6range"
	"github.com/linode/terraform-provider-linode/v2/linode/ipv6ranges"
	"github.com/linode/terraform-provider-linode/v2/linode/kernel"
	"github.com/linode/terraform-provider-linode/v2/linode/kernels"
	"github.com/linode/terraform-provider-linode/v2/linode/lke"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeclusters"
	"github.com/linode/terraform-provider-linode/v2/linode/lkenodepool"
	"github.com/linode/terraform-provider-linode/v2/linode/lketypes"
	"github.com/linode/terraform-provider-linode/v2/linode/lkeversions"
	"github.com/linode/terraform-provider-linode/v2/linode/nb"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfigs"
	"github.com/linode/terraform-provider-linode/v2/linode/nbnode"
	"github.com/linode/terraform-provider-linode/v2/linode/nbs"
	"github.com/linode/terraform-provider-linode/v2/linode/nbtypes"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingip"
	"github.com/linode/terraform-provider-linode/v2/linode/networkingips"
	"github.com/linode/terraform-provider-linode/v2/linode/networkipassignment"
	"github.com/linode/terraform-provider-linode/v2/linode/networkreservedip"
	"github.com/linode/terraform-provider-linode/v2/linode/networkreservedips"
	"github.com/linode/terraform-provider-linode/v2/linode/networktransferprices"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket"
	"github.com/linode/terraform-provider-linode/v2/linode/objcluster"
	"github.com/linode/terraform-provider-linode/v2/linode/objkey"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroup"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroupassignment"
	"github.com/linode/terraform-provider-linode/v2/linode/placementgroups"
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
	"github.com/linode/terraform-provider-linode/v2/linode/volumetypes"
	"github.com/linode/terraform-provider-linode/v2/linode/vpc"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcips"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcs"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnet"
	"github.com/linode/terraform-provider-linode/v2/linode/vpcsubnets"
)

type FrameworkProvider struct {
	ProviderVersion string
	Meta            *helper.FrameworkProviderMeta
}

// CreateFrameworkProviderWithMeta is used by the crossplane provider
func CreateFrameworkProviderWithMeta(version string, meta *helper.ProviderMeta) provider.ProviderWithValidateConfig {
	return &FrameworkProvider{
		ProviderVersion: version,
		Meta: &helper.FrameworkProviderMeta{
			Client: &meta.Client,
			Config: helper.GetFrameworkProviderModelFromSDKv2ProviderConfig(meta.Config),
		},
	}
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
	resp.TypeName = "linode"
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
				Optional:    true,
				Description: "The path to the Linode config file to use. (default `~/.config/linode`)",
			},
			"config_profile": schema.StringAttribute{
				Optional:    true,
				Description: "The Linode config profile to use. (default `default`)",
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
			"obj_access_key": schema.StringAttribute{
				Optional:    true,
				Description: "The access key to be used in linode_object_storage_bucket and linode_object_storage_object.",
			},
			"obj_secret_key": schema.StringAttribute{
				Optional:    true,
				Description: "The secret key to be used in linode_object_storage_bucket and linode_object_storage_object.",
				Sensitive:   true,
			},
			"obj_use_temp_keys": schema.BoolAttribute{
				Optional: true,
				Description: "If true, temporary object keys will be created implicitly at apply-time " +
					"for the linode_object_storage_object and linode_object_sorage_bucket resource.",
			},
			"obj_bucket_force_delete": schema.BoolAttribute{
				Optional: true,
				Description: "If true, when deleting a linode_object_storage_bucket any objects " +
					"and versions will be force deleted.",
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
		volume.NewResource,
		instancesharedips.NewResource,
		instancedisk.NewResource,
		lkenodepool.NewResource,
		image.NewResource,
		nbconfig.NewResource,
		firewall.NewResource,
		placementgroup.NewResource,
		placementgroupassignment.NewResource,
		instancereservedip.NewResource,
		networkreservedip.NewResource,
		networkingip.NewResource,
		networkipassignment.NewResource,
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
		networktransferprices.NewDataSource,
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
		nbtypes.NewDataSource,
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
		vpcips.NewDataSource,
		vpcsubnets.NewDataSource,
		vpcs.NewDataSource,
		volumes.NewDataSource,
		volumetypes.NewDataSource,
		accountavailability.NewDataSource,
		nbconfigs.NewDataSource,
		ipv6ranges.NewDataSource,
		domains.NewDataSource,
		lke.NewDataSource,
		lkeclusters.NewDataSource,
		lketypes.NewDataSource,
		placementgroup.NewDataSource,
		placementgroups.NewDataSource,
		childaccount.NewDataSource,
		childaccounts.NewDataSource,
		networkreservedip.NewDataSourceFetch,
		networkreservedips.NewDataSourceList,
		networkingips.NewDataSource,
	}
}
