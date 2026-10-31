//go:build unit

package instance

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linode/linodego/v2"
)

func newConfigUpdateTestServer(t *testing.T, putBodies *[]map[string]any) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v4/linode/instances/123/configs":
			_, _ = io.WriteString(w, `{"data": [{
				"id": 456,
				"label": "My Config",
				"kernel": "linode/grub2",
				"run_level": "default",
				"virt_mode": "paravirt",
				"root_device": "/dev/sda",
				"helpers": {
					"updatedb_disabled": true,
					"distro": true,
					"modules_dep": true,
					"network": false,
					"devtmpfs_automount": false
				},
				"interfaces": []
			}], "page": 1, "pages": 1, "results": 1}`)
		case r.Method == http.MethodPut && r.URL.Path == "/v4/linode/instances/123/configs/456":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("failed to decode config update body: %v", err)
			}
			*putBodies = append(*putBodies, body)
			_, _ = io.WriteString(w, `{"id": 456, "label": "My Config"}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	return srv
}

func testTFConfig(network bool) map[string]any {
	return map[string]any{
		"label":        "My Config",
		"kernel":       "linode/grub2",
		"run_level":    "default",
		"virt_mode":    "paravirt",
		"root_device":  "/dev/sda",
		"comments":     "",
		"memory_limit": 0,
		"helpers": []any{
			map[string]any{
				"updatedb_disabled":  true,
				"distro":             true,
				"modules_dep":        true,
				"network":            network,
				"devtmpfs_automount": false,
			},
		},
		"interface": []any{},
	}
}

func TestUpdateInstanceConfigsHelpers(t *testing.T) {
	tests := []struct {
		name        string
		network     bool
		wantHelpers bool
	}{
		{name: "unchanged helpers are not sent", network: false, wantHelpers: false},
		{name: "changed helpers are sent", network: true, wantHelpers: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var putBodies []map[string]any
			srv := newConfigUpdateTestServer(t, &putBodies)

			client, err := linodego.NewClient(srv.Client())
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			client.SetBaseURL(srv.URL)

			tfConfigs := []any{testTFConfig(tt.network)}

			_, _, _, err = updateInstanceConfigs(
				context.Background(), client, nil, linodego.Instance{ID: 123},
				[]any{testTFConfig(false)}, tfConfigs, map[string]int{}, "",
			)
			if err != nil {
				t.Fatalf("updateInstanceConfigs returned an error: %v", err)
			}

			if len(putBodies) != 1 {
				t.Fatalf("expected 1 config update request, got %d", len(putBodies))
			}

			helpers, found := putBodies[0]["helpers"]
			if found != tt.wantHelpers {
				t.Fatalf("expected helpers in request body: %v, got %v (body: %v)", tt.wantHelpers, found, putBodies[0])
			}

			if tt.wantHelpers && helpers.(map[string]any)["network"] != true {
				t.Fatalf("expected helpers.network to be true, got %v", helpers)
			}
		})
	}
}
