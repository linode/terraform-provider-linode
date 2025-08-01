package helper

func Must[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}

	return result
}
