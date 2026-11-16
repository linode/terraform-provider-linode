//go:build unit

package acceptance

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckDatabaseDestroyResult(t *testing.T) {
	transportErr := errors.New("connection failed")
	serverErr := &linodego.Error{Code: http.StatusInternalServerError, Message: "server error"}
	notFoundErr := &linodego.Error{Code: http.StatusNotFound, Message: "not found"}

	tests := []struct {
		name        string
		err         error
		wantErr     bool
		wantWrapped error
	}{
		{name: "database still exists", wantErr: true},
		{name: "not found", err: notFoundErr},
		{name: "wrapped not found", err: fmt.Errorf("request failed: %w", notFoundErr)},
		{name: "server error", err: serverErr, wantErr: true, wantWrapped: serverErr},
		{name: "transport error", err: transportErr, wantErr: true, wantWrapped: transportErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkDatabaseDestroyResult(123, test.err)
			if !test.wantErr {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), "123")
			if test.wantWrapped != nil {
				assert.ErrorIs(t, err, test.wantWrapped)
			}
		})
	}
}
