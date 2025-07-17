package cmap

type CmapField[K comparable, V any] interface {
	Key() K
	Value() V
}

type cmapFieldImpl[K comparable, V any] struct {
	k K
	v V
}

func (t *cmapFieldImpl[K, V]) Key() K {
	return t.k
}

func (t *cmapFieldImpl[K, V]) Value() V {
	return t.v
}
