package slices

func Map[S ~[]F, F any, T any](s S, f func(F) T) []T {
	result := make([]T, 0, len(s))
	for _, v := range s {
		result = append(result, f(v))
	}

	return result
}
