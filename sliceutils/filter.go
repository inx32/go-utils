package sliceutils

func Filter[T any](s []T, f func(T) bool) (r []T) {
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}
