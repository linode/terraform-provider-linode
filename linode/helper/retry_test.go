//go:build unit

package helper

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestInstanceDiskCreateBusyRetry(t *testing.T) {
	retryFunc := InstanceDiskCreateBusyRetry()

	tests := []struct {
		name          string
		statusCode    int
		url           string
		body          string
		expectedRetry bool
		description   string
	}{
		{
			name:          "Should retry on Linode busy error during disk creation",
			statusCode:    400,
			url:           "/v4/linode/instances/12345/disks",
			body:          `{"errors": [{"reason": "Linode busy."}]}`,
			expectedRetry: true,
			description:   "400 with 'Linode busy.' message on disk creation endpoint",
		},
		{
			name:          "Should not retry on different 400 error",
			statusCode:    400,
			url:           "/v4/linode/instances/12345/disks",
			body:          `{"errors": [{"reason": "Invalid disk size"}]}`,
			expectedRetry: false,
			description:   "400 with different error message",
		},
		{
			name:          "Should not retry on 500 error",
			statusCode:    500,
			url:           "/v4/linode/instances/12345/disks",
			body:          `{"errors": [{"reason": "Linode busy."}]}`,
			expectedRetry: false,
			description:   "Different status code, even with correct message",
		},
		{
			name:          "Should not retry on different endpoint",
			statusCode:    400,
			url:           "/v4/linode/instances/12345",
			body:          `{"errors": [{"reason": "Linode busy."}]}`,
			expectedRetry: false,
			description:   "Correct status and message but wrong endpoint",
		},
		{
			name:          "Should not retry on disk get endpoint",
			statusCode:    400,
			url:           "/v4/linode/instances/12345/disks/67890",
			body:          `{"errors": [{"reason": "Linode busy."}]}`,
			expectedRetry: false,
			description:   "Disk GET endpoint should not match",
		},
		{
			name:          "Should handle URL with query parameters",
			statusCode:    400,
			url:           "/v4/linode/instances/12345/disks?page=1",
			body:          `{"errors": [{"reason": "Linode busy."}]}`,
			expectedRetry: true,
			description:   "URL with query parameters should still match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			parsedURL, _ := url.Parse(tt.url)
			req := &http.Request{
				URL: parsedURL,
			}

			resp := &http.Response{
				StatusCode: tt.statusCode,
				Request:    req,
				Body:       io.NopCloser(bytes.NewReader([]byte(tt.body))),
			}

			result := retryFunc(resp, nil)

			if result != tt.expectedRetry {
				t.Errorf("%s: expected retry=%v, got retry=%v",
					tt.description, tt.expectedRetry, result)
			}
		})
	}
}

func TestInstanceDiskCreateBusyRetry_NilRequest(t *testing.T) {
	retryFunc := InstanceDiskCreateBusyRetry()

	resp := &http.Response{
		StatusCode: 400,
		Request:    nil,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"errors": [{"reason": "Linode busy."}]}`))),
	}

	result := retryFunc(resp, nil)
	if result {
		t.Error("Should not retry when request is nil")
	}
}

func TestInstanceDiskCreateBusyRetry_EOFErrors(t *testing.T) {
	retryFunc := InstanceDiskCreateBusyRetry()

	tests := []struct {
		name          string
		url           string
		err           error
		expectedRetry bool
		description   string
	}{
		{
			name:          "Should retry on EOF error for disk creation",
			url:           "/v4/linode/instances/12345/disks",
			err:           errors.New("unexpected EOF"),
			expectedRetry: true,
			description:   "EOF error on disk creation endpoint",
		},
		{
			name:          "Should retry on EOF with error code",
			url:           "/v4/linode/instances/12345/disks",
			err:           errors.New("[002] unexpected EOF"),
			expectedRetry: true,
			description:   "EOF error with [002] code on disk creation endpoint",
		},
		{
			name:          "Should not retry EOF on different endpoint",
			url:           "/v4/linode/instances/12345",
			err:           errors.New("unexpected EOF"),
			expectedRetry: false,
			description:   "EOF error on wrong endpoint",
		},
		{
			name:          "Should not retry EOF on disk get endpoint",
			url:           "/v4/linode/instances/12345/disks/67890",
			err:           errors.New("unexpected EOF"),
			expectedRetry: false,
			description:   "EOF error on disk GET endpoint",
		},
		{
			name:          "Should not retry non-EOF errors",
			url:           "/v4/linode/instances/12345/disks",
			err:           errors.New("connection refused"),
			expectedRetry: false,
			description:   "Different error type should not retry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, _ := url.Parse(tt.url)
			req := &http.Request{
				URL: parsedURL,
			}

			resp := &http.Response{
				StatusCode: 0,
				Request:    req,
			}

			result := retryFunc(resp, tt.err)

			if result != tt.expectedRetry {
				t.Errorf("%s: expected retry=%v, got retry=%v",
					tt.description, tt.expectedRetry, result)
			}
		})
	}
}

func TestInstanceDiskCreateBusyRetry_EOFWithNilResponse(t *testing.T) {
	retryFunc := InstanceDiskCreateBusyRetry()

	// EOF error with nil response should not retry
	result := retryFunc(nil, errors.New("unexpected EOF"))
	if result {
		t.Error("Should not retry EOF when response is nil")
	}
}
