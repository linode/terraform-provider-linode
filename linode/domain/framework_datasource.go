package domain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type DataSource struct {
	helper.BaseDataSource
}

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_domain",
				Schema: &frameworkDataSourceSchema,
			},
		),
	}
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_domain")

	var data DomainModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var domain *linodego.Domain
	var err diag.Diagnostic

	// Resolve the domain from the corresponding field
	if !data.ID.IsNull() {
		id := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "domain_id", id)

		domain, err = d.getDomainByID(ctx, id)
	} else {
		domain, err = d.getDomainByDomain(ctx, data.Domain.ValueString())
	}

	if err != nil {
		resp.Diagnostics.Append(err)
		return
	}

	data.ParseDomain(domain)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSource) getDomainByID(ctx context.Context, id int) (*linodego.Domain, diag.Diagnostic) {
	tflog.Trace(ctx, "client.GetDomain(...)")

	domain, err := d.Client.GetDomain(ctx, id)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed to get Domain with id %d", id),
			err.Error(),
		)
	}

	return domain, nil
}

func (d *DataSource) getDomainByDomain(ctx context.Context, domain string) (*linodego.Domain, diag.Diagnostic) {
	tflog.Trace(ctx, "Get domain by domain", map[string]any{
		"domain": domain,
	})

	filter, err := json.Marshal(map[string]interface{}{"domain": domain})
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			"Failed to marshal domain filter",
			err.Error(),
		)
	}

	domains, err := d.Client.ListDomains(ctx, linodego.NewListOptions(0, string(filter)))
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
