package sliceutils

import "slices"

func RemoveIndex[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

func RemoveIndexFast[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveIndexFastC[T any](s []T, i int) []T {
	return RemoveIndexFast(slices.Clone(s), i)
}
