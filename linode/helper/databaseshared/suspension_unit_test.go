//go:build unit

package databaseshared

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
)

func newSuspensionTestClient(t *testing.T, handler http.Handler) *linodego.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := linodego.NewClient(server.Client())
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	client.SetBaseURL(server.URL)
	client.SetPollDelay(time.Millisecond)
	return &client
}

func writeMockSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("{}"))
}

func TestReconcileSuspensionSyncNoop(t *testing.T) {
	for _, suspended := range []bool{false, true} {
		t.Run(map[bool]string{false: "active", true: "suspended"}[suspended], func(t *testing.T) {
			err := ReconcileSuspensionSync(
				context.Background(), nil, 123, linodego.DatabaseEngineTypeValkey, suspended, suspended, time.Second,
			)
			if err != nil {
				t.Fatalf("unexpected error for no-op reconciliation: %v", err)
			}
		})
	}
}

func TestReconcileSuspensionSyncValkeyDispatch(t *testing.T) {
	tests := []struct {
		name           string
		databaseStatus bool
		desiredStatus  bool
		operation      string
		apiStatus      linodego.DatabaseStatus
	}{
		{name: "suspend", desiredStatus: true, operation: "suspend", apiStatus: linodego.DatabaseStatusSuspended},
		{name: "resume", databaseStatus: true, operation: "resume", apiStatus: linodego.DatabaseStatusActive},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			operationCalled := false
			client := newSuspensionTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					operationCalled = true
					if !strings.HasSuffix(r.URL.Path, "/"+test.operation) {
						t.Errorf("unexpected operation path %q", r.URL.Path)
					}
					writeMockSuccess(w)
					return
				}

				if r.Method != http.MethodGet {
					t.Errorf("unexpected request method %q", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/instances/123") {
					t.Errorf("unexpected status path %q", r.URL.Path)
				}
				_ = json.NewEncoder(w).Encode(linodego.ValkeyDatabase{ID: 123, Status: test.apiStatus})
			}))

			err := ReconcileSuspensionSync(
				context.Background(), client, 123, linodego.DatabaseEngineTypeValkey,
				test.databaseStatus, test.desiredStatus, time.Second,
			)
			if err != nil {
				t.Fatalf("unexpected reconciliation error: %v", err)
			}
			if !operationCalled {
				t.Fatal("expected Valkey suspend or resume operation")
			}
		})
	}
}

func TestReconcileSuspensionSyncPropagatesErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		failOnGet  bool
	}{
		{name: "suspend operation", statusCode: http.StatusInternalServerError},
		{name: "status wait", statusCode: http.StatusInternalServerError, failOnGet: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := newSuspensionTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					if test.failOnGet {
						http.Error(w, "status failed", test.statusCode)
						return
					}
					_ = json.NewEncoder(w).Encode(linodego.ValkeyDatabase{ID: 123, Status: linodego.DatabaseStatusSuspended})
					return
				}
				if !test.failOnGet {
					http.Error(w, "suspend failed", test.statusCode)
					return
				}
				writeMockSuccess(w)
			}))

			err := ReconcileSuspensionSync(
				context.Background(), client, 123, linodego.DatabaseEngineTypeValkey, false, true, time.Second,
			)
			if err == nil {
				t.Fatal("expected reconciliation error")
			}
		})
	}
}

func TestReconcileSuspensionSyncPropagatesDeadline(t *testing.T) {
	client := newSuspensionTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			writeMockSuccess(w)
			return
		}
		<-r.Context().Done()
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := ReconcileSuspensionSync(ctx, client, 123, linodego.DatabaseEngineTypeValkey, false, true, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
