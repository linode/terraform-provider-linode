package helper

func TypedSliceToAny[T any](obj []T) []any {
	result := make([]any, len(obj))

	for i, v := range obj {
		result[i] = v
	}

	return result
}
