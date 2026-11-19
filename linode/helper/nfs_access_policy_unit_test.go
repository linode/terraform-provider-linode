//go:build unit

package helper

import (
	"context"
	"testing"
	"time"

	"github.com/linode/linodego/v2"
)

func TestWaitForNFSAccessPolicySettled(t *testing.T) {
	cases := []struct {
		status  linodego.NFSAccessPolicyStatus
		wantErr bool
	}{
		{linodego.NFSAccessPolicyStatusActive, false},
		{"pending", false},
		{"", false},
		{linodego.NFSAccessPolicyStatusError, true},
	}

	for _, c := range cases {
		t.Run(string(c.status), func(t *testing.T) {
			_, err := WaitForNFSAccessPolicySettled(
				context.Background(),
				time.Minute,
				func(context.Context) (*struct{}, linodego.NFSAccessPolicyStatus, error) {
					return &struct{}{}, c.status, nil
				},
			)
			if (err != nil) != c.wantErr {
				t.Fatalf("unexpected error result: %v", err)
			}
		})
	}
}

func TestWaitForNFSAccessPolicySettled_timeout(t *testing.T) {
	_, err := WaitForNFSAccessPolicySettled(
		context.Background(),
		10*time.Millisecond,
		func(context.Context) (*struct{}, linodego.NFSAccessPolicyStatus, error) {
			return &struct{}{}, linodego.NFSAccessPolicyStatusUpdating, nil
		},
	)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestRetryNFSAccessPolicyUpdate(t *testing.T) {
	notReady := &linodego.Error{
		Code:    400,
		Message: "NFS space is not ready for an access policy yet; try again once it is provisioned",
	}

	calls := 0
	err := RetryNFSAccessPolicyUpdate(context.Background(), time.Minute, func(context.Context) error {
		calls++
		if calls < 2 {
			return notReady
		}
		return nil
	})
	if err != nil || calls != 2 {
		t.Fatalf("expected success after 2 calls, got err=%v calls=%d", err, calls)
	}

	calls = 0
	other := &linodego.Error{Code: 400, Message: "invalid squash_policy"}
	if err := RetryNFSAccessPolicyUpdate(context.Background(), time.Minute, func(context.Context) error {
		calls++
		return other
	}); err == nil || calls != 1 {
		t.Fatalf("expected immediate failure, got err=%v calls=%d", err, calls)
	}

	if err := RetryNFSAccessPolicyUpdate(context.Background(), 10*time.Millisecond, func(context.Context) error {
		return notReady
	}); err == nil {
		t.Fatal("expected timeout error")
	}
}
