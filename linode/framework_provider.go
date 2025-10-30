package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/account"
	"github.com/linode/terraform-provider-linode/v3/linode/accountavailabilities"
	"github.com/linode/terraform-provider-linode/v3/linode/accountavailability"
	"github.com/linode/terraform-provider-linode/v3/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/v3/linode/accountlogins"
	"github.com/linode/terraform-provider-linode/v3/linode/accountsettings"
	"github.com/linode/terraform-provider-linode/v3/linode/backup"
	"github.com/linode/terraform-provider-linode/v3/linode/childaccount"
	"github.com/linode/terraform-provider-linode/v3/linode/childaccounts"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegroup"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegroupimageshares"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegrouptoken"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegrouptokens"
	"github.com/linode/terraform-provider-linode/v3/linode/databasebackups"
	"github.com/linode/terraform-provider-linode/v3/linode/databaseengines"
	"github.com/linode/terraform-provider-linode/v3/linode/databasemysql"
	"github.com/linode/terraform-provider-linode/v3/linode/databasemysqlconfig"
	"github.com/linode/terraform-provider-linode/v3/linode/databasemysqlv2"
	"github.com/linode/terraform-provider-linode/v3/linode/databasepostgresql"
	"github.com/linode/terraform-provider-linode/v3/linode/databasepostgresqlconfig"
	"github.com/linode/terraform-provider-linode/v3/linode/databasepostgresqlv2"
	"github.com/linode/terraform-provider-linode/v3/linode/databases"
	"github.com/linode/terraform-provider-linode/v3/linode/domain"
	"github.com/linode/terraform-provider-linode/v3/linode/domainrecord"
	"github.com/linode/terraform-provider-linode/v3/linode/domains"
	"github.com/linode/terraform-provider-linode/v3/linode/domainzonefile"
	"github.com/linode/terraform-provider-linode/v3/linode/firewall"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalldevice"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalls"
	"github.com/linode/terraform-provider-linode/v3/linode/firewallsettings"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalltemplate"
	"github.com/linode/terraform-provider-linode/v3/linode/firewalltemplates"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
	"github.com/linode/terraform-provider-linode/v3/linode/image"
	"github.com/linode/terraform-provider-linode/v3/linode/images"
	"github.com/linode/terraform-provider-linode/v3/linode/instancedisk"
	"github.com/linode/terraform-provider-linode/v3/linode/instanceip"
	"github.com/linode/terraform-provider-linode/v3/linode/instancenetworking"
	"github.com/linode/terraform-provider-linode/v3/linode/instancereservedipassignment"
	"github.com/linode/terraform-provider-linode/v3/linode/instancesharedips"
	"github.com/linode/terraform-provider-linode/v3/linode/instancetype"
	"github.com/linode/terraform-provider-linode/v3/linode/instancetypes"
	"github.com/linode/terraform-provider-linode/v3/linode/ipv6range"
	"github.com/linode/terraform-provider-linode/v3/linode/ipv6ranges"
	"github.com/linode/terraform-provider-linode/v3/linode/kernel"
	"github.com/linode/terraform-provider-linode/v3/linode/kernels"
	"github.com/linode/terraform-provider-linode/v3/linode/lke"
	"github.com/linode/terraform-provider-linode/v3/linode/lkeclusters"
	"github.com/linode/terraform-provider-linode/v3/linode/lkenodepool"
	"github.com/linode/terraform-provider-linode/v3/linode/lketypes"
	"github.com/linode/terraform-provider-linode/v3/linode/lkeversion"
	"github.com/linode/terraform-provider-linode/v3/linode/lkeversions"
	"github.com/linode/terraform-provider-linode/v3/linode/maintenancepolicies"
	"github.com/linode/terraform-provider-linode/v3/linode/nb"
	"github.com/linode/terraform-provider-linode/v3/linode/nbconfig"
	"github.com/linode/terraform-provider-linode/v3/linode/nbconfigs"
	"github.com/linode/terraform-provider-linode/v3/linode/nbnode"
	"github.com/linode/terraform-provider-linode/v3/linode/nbs"
	"github.com/linode/terraform-provider-linode/v3/linode/nbtypes"
	"github.com/linode/terraform-provider-linode/v3/linode/nbvpc"
	"github.com/linode/terraform-provider-linode/v3/linode/nbvpcs"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingip"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingipassignment"
	"github.com/linode/terraform-provider-linode/v3/linode/networkingips"
	"github.com/linode/terraform-provider-linode/v3/linode/networktransferprices"
	"github.com/linode/terraform-provider-linode/v3/linode/obj"
	"github.com/linode/terraform-provider-linode/v3/linode/objbucket"
	"github.com/linode/terraform-provider-linode/v3/linode/objcluster"
	"github.com/linode/terraform-provider-linode/v3/linode/objendpoints"
	"github.com/linode/terraform-provider-linode/v3/linode/objkey"
	"github.com/linode/terraform-provider-linode/v3/linode/objquota"
	"github.com/linode/terraform-provider-linode/v3/linode/objquotas"
	"github.com/linode/terraform-provider-linode/v3/linode/placementgroup"
	"github.com/linode/terraform-provider-linode/v3/linode/placementgroupassignment"
	"github.com/linode/terraform-provider-linode/v3/linode/placementgroups"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroup"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupimageshares"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmember"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmembers"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroups"
	"github.com/linode/terraform-provider-linode/v3/linode/profile"
	"github.com/linode/terraform-provider-linode/v3/linode/rdns"
	"github.com/linode/terraform-provider-linode/v3/linode/region"
	"github.com/linode/terraform-provider-linode/v3/linode/regions"
	"github.com/linode/terraform-provider-linode/v3/linode/sshkey"
	"github.com/linode/terraform-provider-linode/v3/linode/sshkeys"
	"github.com/linode/terraform-provider-linode/v3/linode/stackscript"
	"github.com/linode/terraform-provider-linode/v3/linode/stackscripts"
	"github.com/linode/terraform-provider-linode/v3/linode/token"
	"github.com/linode/terraform-provider-linode/v3/linode/user"
	"github.com/linode/terraform-provider-linode/v3/linode/users"
	"github.com/linode/terraform-provider-linode/v3/linode/vlan"
	"github.com/linode/terraform-provider-linode/v3/linode/volume"
	"github.com/linode/terraform-provider-linode/v3/linode/volumes"
	"github.com/linode/terraform-provider-linode/v3/linode/volumetypes"
	"github.com/linode/terraform-provider-linode/v3/linode/vpc"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcips"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcs"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcsubnet"
	"github.com/linode/terraform-provider-linode/v3/linode/vpcsubnets"
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
	resp.Version = p.ProviderVersion
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
			"api_ca_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path to a Linode API CA file to trust.",
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
		accountsettings.NewResource,
		firewall.NewResource,
		firewalldevice.NewResource,
		image.NewResource,
		instancedisk.NewResource,
		instanceip.NewResource,
		instancesharedips.NewResource,
		ipv6range.NewResource,
		lkenodepool.NewResource,
		nb.NewResource,
		nbconfig.NewResource,
		nbnode.NewResource,
		objkey.NewResource,
		placementgroup.NewResource,
		placementgroupassignment.NewResource,
		instancereservedipassignment.NewResource,
		rdns.NewResource,
		sshkey.NewResource,
		stackscript.NewResource,
		token.NewResource,
		volume.NewResource,
		vpc.NewResource,
		vpcsubnet.NewResource,
		databasepostgresqlv2.NewResource,
		networkingip.NewResource,
		networkingipassignment.NewResource,
		obj.NewResource,
		databasemysqlv2.NewResource,
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
		lkeversion.NewDataSource,
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
		nbvpc.NewDataSource,
		nbvpcs.NewDataSource,
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
		maintenancepolicies.NewDataSource,
		placementgroup.NewDataSource,
		placementgroups.NewDataSource,
		childaccount.NewDataSource,
		childaccounts.NewDataSource,
		networkingips.NewDataSource,
		databasemysqlv2.NewDataSource,
		databasepostgresqlv2.NewDataSource,
		databasemysqlconfig.NewDataSource,
		databasepostgresqlconfig.NewDataSource,
		objendpoints.NewDataSource,
		objquota.NewDataSource,
		objquotas.NewDataSource,
		firewalltemplate.NewDataSource,
		firewalltemplates.NewDataSource,
		firewallsettings.NewDataSource,
		producerimagesharegroup.NewDataSource,
		producerimagesharegroups.NewDataSource,
		producerimagesharegroupimageshares.NewDataSource,
		producerimagesharegroupmember.NewDataSource,
		producerimagesharegroupmembers.NewDataSource,
		consumerimagesharegrouptoken.NewDataSource,
		consumerimagesharegrouptokens.NewDataSource,
		consumerimagesharegroup.NewDataSource,
		consumerimagesharegroupimageshares.NewDataSource,
	}
}
