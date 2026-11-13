//go:build unit

package consumerimagesharegrouptoken

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/require"
)

func TestReadRemovesMissingImageShareGroupToken(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors": [{"reason": "Not found"}]}`))
	}))
	t.Cleanup(server.Close)

	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetRetryCount(0)

	objType := frameworkResourceSchema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{
		"token_uuid": tftypes.NewValue(tftypes.String, "5c7a3b6e-1f2d-4e8a-9b0c-123456789abc"),
	}
	for name, typ := range objType.AttributeTypes {
		if _, ok := values[name]; !ok {
			values[name] = tftypes.NewValue(typ, nil)
		}
	}
	state := tfsdk.State{
		Schema: frameworkResourceSchema,
		Raw:    tftypes.NewValue(objType, values),
	}

	r := NewResource().(*Resource)
	r.Meta = &helper.FrameworkProviderMeta{Client: &client}

	resp := resource.ReadResponse{State: state}
	r.Read(ctx, resource.ReadRequest{State: state}, &resp)

	require.False(t, resp.Diagnostics.HasError(), "unexpected error: %v", resp.Diagnostics)
	require.True(t, resp.State.Raw.IsNull(), "expected the resource to be removed from state")
}
