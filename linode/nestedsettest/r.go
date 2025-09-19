package nestedsettest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknownIf returns a plan modifier that will only use the state value
// in place of an unknown value if the given condition is met.
func SampleModifier() planmodifier.String {
	return sampleModifier{}
}

type sampleModifier struct{}

func (m sampleModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change " +
		"if the given condition is met."
}

func (m sampleModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change " +
		"if the given condition is met."
}

func (m sampleModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	panic(fmt.Sprintf("%v", req.PathExpression))
}

func NewResource() resource.Resource {
	return &setNestedExampleResource{}
}

type setNestedExampleResource struct{}

func (r *setNestedExampleResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = "linode_nested_attrs_resource"
}

func (r *setNestedExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nested": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"double_nested": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"str_a": schema.StringAttribute{
								Computed: true,
								Optional: true,
								PlanModifiers: []planmodifier.String{
									SampleModifier(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *setNestedExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *setNestedExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.State.Raw = req.State.Raw
}

func (r *setNestedExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.State.Raw = req.Plan.Raw
}

func (r *setNestedExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
