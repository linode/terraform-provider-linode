//go:build integration

package lke_test

//
//func TestReconcileLKENodePoolSpecs(t *testing.T) {
//	for _, tc := range []struct {
//		name             string
//		specs            []lke.NodePoolSpec
//		provisionedPools []linodego.LKENodePool
//
//		expectedToDelete []int
//		expectedToCreate []linodego.LKENodePoolCreateOptions
//		expectedToUpdate map[int]linodego.LKENodePoolUpdateOptions
//	}{
//		{
//			name: "no change",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-1", Count: 2},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-1", Count: 2},
//			},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{},
//		},
//		{
//			name: "upsize a single pool",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-1", Count: 2},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-1", Count: 3},
//			},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
//				123: {Count: 3},
//			},
//		},
//		{
//			name: "change single pool type",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-1", Count: 2},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-2", Count: 2},
//			},
//			expectedToCreate: []linodego.LKENodePoolCreateOptions{
//				{Type: "g6-standard-2", Count: 2},
//			},
//			expectedToDelete: []int{123},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{},
//		},
//		{
//			name: "reuse cluster for resize",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-1", Count: 1},
//				{ID: 124, Type: "g6-standard-1", Count: 10},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-1", Count: 9},  // bumped from 1 to 9
//				{Type: "g6-standard-2", Count: 10}, // type changed
//			},
//			expectedToDelete: []int{123},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
//				124: {Count: 9},
//			},
//			expectedToCreate: []linodego.LKENodePoolCreateOptions{
//				{Type: "g6-standard-2", Count: 10},
//			},
//		},
//		{
//			name: "competing resizes",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-3", Count: 3},
//				{ID: 124, Type: "g6-standard-3", Count: 7},
//				{ID: 126, Type: "g6-standard-3", Count: 4},
//				{ID: 127, Type: "g6-standard-3", Count: 2},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-3", Count: 2},
//				{Type: "g6-standard-3", Count: 9},
//				{Type: "g6-standard-3", Count: 8},
//				{Type: "g6-standard-3", Count: 2},
//			},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
//				123: {Count: 2}, // -1
//				124: {Count: 8}, // +1
//				126: {Count: 9}, // +5
//			},
//		},
//		{
//			name: "scaler",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-3", Count: 3},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-3", Count: 3, AutoScalerEnabled: true, AutoScalerMin: 3, AutoScalerMax: 7},
//			},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
//				123: {Count: 3, Autoscaler: &linodego.LKENodePoolAutoscaler{Enabled: true, Min: 3, Max: 7}}, // -1
//			},
//		},
//		{
//			name: "scaler drop",
//			provisionedPools: []linodego.LKENodePool{
//				{ID: 123, Type: "g6-standard-3", Count: 3, Autoscaler: linodego.LKENodePoolAutoscaler{Enabled: true, Min: 3, Max: 7}},
//			},
//			specs: []lke.NodePoolSpec{
//				{Type: "g6-standard-3", Count: 3, AutoScalerEnabled: false},
//			},
//			expectedToUpdate: map[int]linodego.LKENodePoolUpdateOptions{
//				123: {Count: 3, Autoscaler: &linodego.LKENodePoolAutoscaler{Enabled: false, Min: 3, Max: 3}}, // -1
//			},
//		},
//	} {
//		t.Run(tc.name, func(t *testing.T) {
//			updates := lke.ReconcileLKENodePoolSpecs(tc.specs, tc.provisionedPools)
//			if !reflect.DeepEqual(tc.expectedToCreate, updates.ToCreate) {
//				t.Errorf("expected to create:\n%#v\ngot:\n%#v", tc.expectedToCreate, updates.ToCreate)
//			}
//			if !reflect.DeepEqual(tc.expectedToUpdate, updates.ToUpdate) {
//				t.Errorf("expected to update:\n%#v\ngot:\n%#v", tc.expectedToUpdate, updates.ToUpdate)
//			}
//			if !reflect.DeepEqual(tc.expectedToDelete, updates.ToDelete) {
//				t.Errorf("expected to delete:\n%#v\ngot:\n%#v", tc.expectedToDelete, updates.ToDelete)
//			}
//		})
//	}
//}
