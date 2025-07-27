package sliceutils

func Map[T, R any](s []T, f func(T) R) []R {
	if len(s) == 0 {
		return []R{}
	}
	r := make([]R, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}
