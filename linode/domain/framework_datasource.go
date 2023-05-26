package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type DataSource struct {
	client *linodego.Client
}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

func (d *DataSource) Configure(
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

	d.client = meta.Client
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_domain"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data DomainModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domain *linodego.Domain
	var err diag.Diagnostic

	// Resolve the domain from the corresponding field
	if !data.ID.IsNull() {
		domain, err = d.getDomainByID(ctx, int(data.ID.ValueInt64()))
	} else {
		domain, err = d.getDomainByDomain(ctx, data.Domain.ValueString())
	}

	if err != nil {
		resp.Diagnostics.Append(err)
		return
	}

	data.parseDomain(domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSource) getDomainByID(ctx context.Context, id int) (*linodego.Domain, diag.Diagnostic) {
	domain, err := d.client.GetDomain(ctx, id)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed to get Domain with id %d", id),
			err.Error(),
		)
	}

	return domain, nil
}

func (d *DataSource) getDomainByDomain(ctx context.Context, domain string) (*linodego.Domain, diag.Diagnostic) {
	filter, _ := json.Marshal(map[string]interface{}{"domain": domain})
	domains, err := d.client.ListDomains(ctx, linodego.NewListOptions(0, string(filter)))
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			"Failed to list matching domains",
			err.Error(),
		)
	}
	if len(domains) != 1 || domains[0].Domain != domain {
		return nil, diag.NewErrorDiagnostic(
			"Failed to retrieve Linode Domain",
			fmt.Sprintf("Domain %s was not found in list result", domain),
		)
	}
	return &domains[0], nil
}
