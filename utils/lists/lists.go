package lists

func Map[T, E any](list []T, f func(T) E) []E {
	result := make([]E, 0, len(list))

	for _, item := range list {
		result = append(result, f(item))
	}

	return result
}
