package sliceutils

func Map[T, R any](s []T, f func(T) R) (r []R) {
	for _, i := range s {
		r = append(r, f(i))
	}
	return
}
