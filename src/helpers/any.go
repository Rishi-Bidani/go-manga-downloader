package helpers

func Any[T any](ts []T, fn func(T) bool) bool {
	for _, t := range ts {
		if fn(t) {
			return true
		}
	}
	return false
}