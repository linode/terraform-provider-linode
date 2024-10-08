//go:build unit

package lke_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/lke"
)

func TestReconcileLKENodePoolSpecs(t *testing.T) {
	for _, tc := range []struct {
		name     string
		oldSpecs []lke.NodePoolSpec
		newSpecs []lke.NodePoolSpec

		expectedToDelete []int
		expectedToCreate []linodego.LKENodePoolCreateOptions
		expectedToUpdate map[int]linodego.LKENodePoolUpdateOptions
	}{
		{
			name: "no change",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
			expectedToDelete: []int{},
		},
		{
			name: "upsize a single pool",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 3, Tags: []string{"example"}},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 3, Tags: &[]string{"example"}},
			},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
			expectedToDelete: []int{},
		},
		{
			name: "change single pool type",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 2},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-2", Count: 2},
			},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{
				{Type: "g6-standard-2", Count: 2},
			},
			expectedToDelete: []int{123},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{},
		},
		{
			name: "reuse cluster for resize",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 1},
				{ID: 124, Type: "g6-standard-1", Count: 10},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-1", Count: 9, Tags: []string{"example"}},  // bumped from 1 to 9
				{ID: 124, Type: "g6-standard-2", Count: 10, Tags: []string{"example"}}, // type changed
			},
			expectedToDelete: []int{124},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 9, Tags: &[]string{"example"}},
			},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{
				{Type: "g6-standard-2", Count: 10, Tags: []string{"example"}},
			},
		},
		{
			name: "competing resizes",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3},
				{ID: 124, Type: "g6-standard-3", Count: 7},
				{ID: 126, Type: "g6-standard-3", Count: 4},
				{ID: 127, Type: "g6-standard-3", Count: 2},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 2, Tags: []string{"example"}},
				{ID: 124, Type: "g6-standard-3", Count: 9, Tags: []string{"example"}},
				{ID: 126, Type: "g6-standard-3", Count: 8, Tags: []string{"example"}},
				{ID: 127, Type: "g6-standard-3", Count: 2, Tags: []string{"example"}},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 2, Tags: &[]string{"example"}},
				124: {Count: 9, Tags: &[]string{"example"}},
				126: {Count: 8, Tags: &[]string{"example"}},
				127: {Count: 2, Tags: &[]string{"example"}},
			},
			expectedToDelete: []int{},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
		},
		{
			name: "scaler",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3, AutoScalerEnabled: true, AutoScalerMin: 3, AutoScalerMax: 7, Tags: []string{"example"}},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 3, Autoscaler: &linodego.LKENodePoolAutoscaler{Enabled: true, Min: 3, Max: 7}, Tags: &[]string{"example"}},
			},
			expectedToDelete: []int{},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
		},
		{
			name: "scaler drop",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3, AutoScalerEnabled: true, AutoScalerMin: 3, AutoScalerMax: 7},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3, AutoScalerEnabled: false, Tags: []string{"example"}},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 3, Autoscaler: &linodego.LKENodePoolAutoscaler{Enabled: false, Min: 0, Max: 0}, Tags: &[]string{"example"}},
			},
			expectedToDelete: []int{},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
		},
		{
			name: "scaler update",
			oldSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3, AutoScalerEnabled: true, AutoScalerMin: 3, AutoScalerMax: 7},
			},
			newSpecs: []lke.NodePoolSpec{
				{ID: 123, Type: "g6-standard-3", Count: 3, AutoScalerEnabled: true, AutoScalerMin: 5, AutoScalerMax: 10, Tags: []string{"example"}},
			},
			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
				123: {Count: 3, Autoscaler: &linodego.LKENodePoolAutoscaler{Enabled: true, Min: 5, Max: 10}, Tags: &[]string{"example"}},
			},
			expectedToDelete: []int{},
			expectedToCreate: []linodego.LKENodePoolCreateOptions{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			updates, err := lke.ReconcileLKENodePoolSpecs(context.Background(), tc.oldSpecs, tc.newSpecs)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(tc.expectedToCreate, updates.ToCreate) {
				t.Errorf("expected to create:\n%#v\ngot:\n%#v", tc.expectedToCreate, updates.ToCreate)
			}
			if !reflect.DeepEqual(tc.expectedToUpdate, updates.ToUpdate) {
				t.Errorf("expected to update:\n%#v\ngot:\n%#v", tc.expectedToUpdate, updates.ToUpdate)
			}
			if !reflect.DeepEqual(tc.expectedToDelete, updates.ToDelete) {
				t.Errorf("expected to delete:\n%#v\ngot:\n%#v", tc.expectedToDelete, updates.ToDelete)
			}
		})
	}
}
