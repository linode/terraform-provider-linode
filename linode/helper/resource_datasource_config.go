package helper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func GetMetaFromProviderData(
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) *FrameworkProviderMeta {
	meta, ok := req.ProviderData.(*FrameworkProviderMeta)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return nil
	}

	return meta
}

func GetMetaFromProviderDataDatasource(
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) *FrameworkProviderMeta {
	meta, ok := req.ProviderData.(*FrameworkProviderMeta)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return nil
	}

	return meta
}
