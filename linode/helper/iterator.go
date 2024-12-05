package helper

import (
	"iter"
	"maps"
	"slices"
)

// Map returns a new iterator of the values in the given iterator transformed using the given transform function.
func Map[I, O any](values iter.Seq[I], transform func(I) O) iter.Seq[O] {
	return func(yield func(O) bool) {
		for value := range values {
			if !yield(transform(value)) {
				return
			}
		}
	}
}

// Map2 returns a new two-value iterator of the values in the given iterator transformed using the given transform function.
func Map2[I1, I2, O1, O2 any](values iter.Seq2[I1, I2], transform func(I1, I2) (O1, O2)) iter.Seq2[O1, O2] {
	return func(yield func(O1, O2) bool) {
		for value1, value2 := range values {
			if !yield(transform(value1, value2)) {
				return
			}
		}
	}
}

// MapSlice returns a new slice of the values in the given slice transformed using the given transform function.
func MapSlice[I, O any](values []I, transform func(I) O) []O {
	return slices.Collect(Map(slices.Values(values), transform))
}

// MapMap returns a new map of the keys and values from the given map transformed using the given transform function.
func MapMap[IK, OK comparable, IV, OV any](values map[IK]IV, transform func(IK, IV) (OK, OV)) map[OK]OV {
	return maps.Collect(
		Map2(
			maps.All(values),
			func(key IK, value IV) (OK, OV) {
				return transform(key, value)
			},
		),
	)
}
