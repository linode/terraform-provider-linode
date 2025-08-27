package helper

import "log"

func Must[T any](result T, err error) T {
	if err != nil {
		log.Fatalf("helper.Must failed: %s", err)
	}

	return result
}
