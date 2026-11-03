//go:build unit

package accounttfaoptionalusers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testResource(t *testing.T, handler http.HandlerFunc) *Resource {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client, err := linodego.NewClient(server.Client())
	require.NoError(t, err)
	client.SetBaseURL(server.URL)
	client.SetAPIVersion("v4")
	client.SetRetryCount(0)
	r := NewResource().(*Resource)
	r.Meta = &helper.FrameworkProviderMeta{Client: &client}
	return r
}

func testModel(t *testing.T, usernames ...string) ResourceModel {
	t.Helper()
	if usernames == nil {
		usernames = []string{}
	}
	value, diags := types.SetValueFrom(context.Background(), types.StringType, usernames)
	require.False(t, diags.HasError(), "%v", diags)
	return ResourceModel{
		ID:        types.StringValue(singletonID),
		Usernames: value,
	}
}

func testState(t *testing.T, model ResourceModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: frameworkResourceSchema}
	diags := state.Set(context.Background(), &model)
	require.False(t, diags.HasError(), "%v", diags)
	return state
}

func testPlan(t *testing.T, model ResourceModel) tfsdk.Plan {
	t.Helper()
	state := testState(t, model)
	return tfsdk.Plan(state)
}

func assertModel(t *testing.T, state tfsdk.State, expected ResourceModel) {
	t.Helper()
	var actual ResourceModel
	diags := state.Get(context.Background(), &actual)
	require.False(t, diags.HasError(), "%v", diags)
	assert.Equal(t, expected.ID, actual.ID)
	assert.True(t, expected.Usernames.Equal(actual.Usernames), "expected %s, got %s", expected.Usernames, actual.Usernames)
}

func TestResourceLifecycle(t *testing.T) {
	ctx := context.Background()
	users := []string{"previously-exempt"}
	var puts, gets int
	r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/v4/account/tfa-optional-users", req.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch req.Method {
		case http.MethodPut:
			puts++
			var body linodego.TFAOptionalUsersUpdateOptions
			assert.NoError(t, json.NewDecoder(req.Body).Decode(&body))
			assert.NotNil(t, body.Usernames, "an empty set must serialize as [], not null")
			users = slices.Clone(body.Usernames)
			slices.Reverse(users)
			assert.NoError(t, json.NewEncoder(w).Encode(linodego.TFAOptionalUsersUpdateResponse{
				TFAOptionalUsersCount: len(users), URL: "/account/tfa-optional-users",
			}))
		case http.MethodGet:
			gets++
			page, _ := strconv.Atoi(req.URL.Query().Get("page"))
			if page < 1 {
				page = 1
			}
			data := []linodego.TFAOptionalUser{}
			if len(users) > 0 {
				data = append(data, linodego.TFAOptionalUser{Username: users[page-1]})
			}
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"data": data, "page": page, "pages": max(1, len(users)), "results": len(users),
			}))
		default:
			t.Errorf("unexpected method %s", req.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	plan := testModel(t, "alice", "bob")
	plan.ID = types.StringUnknown()
	created := resource.CreateResponse{State: tfsdk.State{Schema: frameworkResourceSchema}}
	r.Create(ctx, resource.CreateRequest{Plan: testPlan(t, plan)}, &created)
	require.False(t, created.Diagnostics.HasError(), "%v", created.Diagnostics)
	assertModel(t, created.State, testModel(t, "alice", "bob"))
	assert.Equal(t, 1, puts)
	assert.Equal(t, 2, gets, "all pages must be read after the count-only PUT response")
	assert.NotContains(t, users, "previously-exempt", "PUT must replace rather than merge")

	users = []string{"bob", "external-user"}
	read := resource.ReadResponse{State: created.State}
	r.Read(ctx, resource.ReadRequest{State: created.State}, &read)
	require.False(t, read.Diagnostics.HasError(), "%v", read.Diagnostics)
	assertModel(t, read.State, testModel(t, "bob", "external-user"))
	assert.Equal(t, 1, puts, "refresh must not write")

	updated := resource.UpdateResponse{State: read.State}
	r.Update(ctx, resource.UpdateRequest{
		Plan: testPlan(t, testModel(t, "alice")), State: read.State,
	}, &updated)
	require.False(t, updated.Diagnostics.HasError(), "%v", updated.Diagnostics)
	assertModel(t, updated.State, testModel(t, "alice"))
	assert.Equal(t, 2, puts)

	empty := resource.UpdateResponse{State: updated.State}
	r.Update(ctx, resource.UpdateRequest{Plan: testPlan(t, testModel(t)), State: updated.State}, &empty)
	require.False(t, empty.Diagnostics.HasError(), "%v", empty.Diagnostics)
	assertModel(t, empty.State, testModel(t))
	assert.Equal(t, 3, puts)

	emptyRead := resource.ReadResponse{State: empty.State}
	r.Read(ctx, resource.ReadRequest{State: empty.State}, &emptyRead)
	require.False(t, emptyRead.Diagnostics.HasError(), "%v", emptyRead.Diagnostics)
	assertModel(t, emptyRead.State, testModel(t))

	users = []string{"imported-user"}
	imported := resource.ImportStateResponse{State: tfsdk.State{Schema: frameworkResourceSchema}}
	r.ImportState(ctx, resource.ImportStateRequest{ID: singletonID}, &imported)
	require.False(t, imported.Diagnostics.HasError(), "%v", imported.Diagnostics)
	importRead := resource.ReadResponse{State: imported.State}
	r.Read(ctx, resource.ReadRequest{State: imported.State}, &importRead)
	require.False(t, importRead.Diagnostics.HasError(), "%v", importRead.Diagnostics)
	assertModel(t, importRead.State, testModel(t, "imported-user"))
	assert.Equal(t, 3, puts, "import must not change exemptions")

	previousPuts, previousGets := puts, gets
	deleted := resource.DeleteResponse{State: importRead.State}
	r.Delete(ctx, resource.DeleteRequest{State: importRead.State}, &deleted)
	require.False(t, deleted.Diagnostics.HasError(), "%v", deleted.Diagnostics)
	assert.Equal(t, []string{"imported-user"}, users)
	assert.Equal(t, previousPuts, puts, "delete must not write exemptions")
	assert.Equal(t, previousGets, gets, "delete must not read exemptions")
}

func TestResourceWriteFailures(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		for _, failMethod := range []string{http.MethodPut, http.MethodGet} {
			t.Run(operation+"/"+failMethod, func(t *testing.T) {
				var methods []string
				r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
					methods = append(methods, req.Method)
					w.Header().Set("Content-Type", "application/json")
					if req.Method == failMethod {
						w.WriteHeader(http.StatusForbidden)
						_, err := w.Write([]byte(`{"errors":[{"reason":"permission denied"}]}`))
						assert.NoError(t, err)
						return
					}
					_, err := w.Write([]byte(`{"tfa_optional_users_count":1,"url":"/account/tfa-optional-users"}`))
					assert.NoError(t, err)
				})
				ctx := context.Background()
				model := testModel(t, "alice")
				var result tfsdk.State
				if operation == "create" {
					resp := resource.CreateResponse{State: tfsdk.State{Schema: frameworkResourceSchema}}
					r.Create(ctx, resource.CreateRequest{Plan: testPlan(t, model)}, &resp)
					require.True(t, resp.Diagnostics.HasError())
					assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "permission denied")
					result = resp.State
				} else {
					previous := testState(t, testModel(t, "previous"))
					resp := resource.UpdateResponse{State: previous}
					r.Update(ctx, resource.UpdateRequest{Plan: testPlan(t, model), State: previous}, &resp)
					require.True(t, resp.Diagnostics.HasError())
					result = resp.State
					if failMethod == http.MethodPut {
						assertModel(t, result, testModel(t, "previous"))
					}
				}
				if failMethod == http.MethodGet {
					assertModel(t, result, model)
					assert.Equal(t, []string{http.MethodPut, http.MethodGet}, methods)
				} else {
					assert.Equal(t, []string{http.MethodPut}, methods)
				}
			})
		}
	}
}

func TestResourceWriteVerifiesMembership(t *testing.T) {
	r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := `{"tfa_optional_users_count":1}`
		if req.Method == http.MethodGet {
			body = `{"data":[{"username":"unexpected-user"}],"page":1,"pages":1,"results":1}`
		}
		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	})
	resp := resource.CreateResponse{State: tfsdk.State{Schema: frameworkResourceSchema}}
	r.Create(context.Background(), resource.CreateRequest{Plan: testPlan(t, testModel(t, "alice"))}, &resp)
	require.True(t, resp.Diagnostics.HasError(), "equal counts do not guarantee equal membership")
	assertModel(t, resp.State, testModel(t, "unexpected-user"))
}

func TestResourceReadFailureRetainsState(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusForbidden} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, err := w.Write([]byte(`{"errors":[{"reason":"unavailable"}]}`))
				assert.NoError(t, err)
			})
			state := testState(t, testModel(t, "alice"))
			resp := resource.ReadResponse{State: state}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			require.True(t, resp.Diagnostics.HasError())
			assertModel(t, resp.State, testModel(t, "alice"))
		})
	}
}

func TestResourceDeleteDoesNotCallAPI(t *testing.T) {
	assert.NotContains(t, frameworkResourceSchema.Attributes, "clear_on_destroy")

	for _, tt := range []struct {
		name      string
		usernames []string
	}{
		{"empty", []string{}},
		{"nonempty", []string{"alice"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var requests atomic.Int32
			r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
				requests.Add(1)
				w.WriteHeader(http.StatusForbidden)
			})
			state := testState(t, testModel(t, tt.usernames...))
			resp := resource.DeleteResponse{State: state}
			r.Delete(context.Background(), resource.DeleteRequest{State: state}, &resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			assert.Zero(t, requests.Load(), "delete must not make any API requests")
		})
	}
}

func TestResourceImportInvalidID(t *testing.T) {
	resp := resource.ImportStateResponse{State: tfsdk.State{Schema: frameworkResourceSchema}}
	NewResource().(*Resource).ImportState(context.Background(), resource.ImportStateRequest{ID: "invalid"}, &resp)
	require.True(t, resp.Diagnostics.HasError())
}

func TestDataSourceRead(t *testing.T) {
	for _, tt := range []struct {
		name      string
		body      string
		usernames []string
	}{
		{"empty", `{"data":[],"page":1,"pages":1,"results":0}`, []string{}},
		{"nonempty", `{"data":[{"username":"alice","url":"/account/users/alice"}],"page":1,"pages":1,"results":1}`, []string{"alice"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := testResource(t, func(w http.ResponseWriter, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write([]byte(tt.body))
				assert.NoError(t, err)
			})
			d := NewDataSource().(*DataSource)
			d.Meta = r.Meta
			resp := datasource.ReadResponse{State: tfsdk.State{Schema: frameworkDataSourceSchema}}
			d.Read(context.Background(), datasource.ReadRequest{}, &resp)
			require.False(t, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
			var model DataSourceModel
			require.False(t, resp.State.Get(context.Background(), &model).HasError())
			assert.Equal(t, types.StringValue(singletonID), model.ID)
			assert.False(t, model.Usernames.IsNull())
			assert.False(t, model.Usernames.IsUnknown())
			assert.True(t, model.Usernames.Equal(testModel(t, tt.usernames...).Usernames))
		})
	}
}
