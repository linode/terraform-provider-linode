package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (r *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client = meta.Client
}

func (r *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_domain"
}

func (r *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
}

func (r *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := r.client

	var data DomainModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domain *linodego.Domain
	var err error

	if !data.ID.IsNull() {
		domain, err = getDomainByID(ctx, client, helper.StringToInt64(data.ID.ValueString(), resp.Diagnostics))
	}

	if !data.Domain.IsNull() {
		domain, err = getDomainByName(ctx, client, data.Domain.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get Domain",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseDomain(ctx, domain)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getDomainByID(ctx context.Context, client *linodego.Client, id int64) (*linodego.Domain, error) {
	domain, err := client.GetDomain(ctx, int(id))
	if err != nil {
		return nil, err
	}

	return domain, nil
}

func getDomainByName(ctx context.Context, client *linodego.Client, domain string) (*linodego.Domain, error) {
	filter, _ := json.Marshal(map[string]interface{}{"domain": domain})

	domains, err := client.ListDomains(ctx, linodego.NewListOptions(0, string(filter)))
	if err != nil {
		return nil, fmt.Errorf("failed to list Domains: %s", err)
	}

	if len(domains) != 1 || domains[0].Domain != domain {
		return nil, fmt.Errorf("domain %s was not found", domain)
	}

	return &domains[0], nil
}
