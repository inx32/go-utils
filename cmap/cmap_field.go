package cmap

type CmapField[K comparable, V any] struct {
	k K
	v V
}

func (t *CmapField[K, V]) Key() K {
	return t.k
}

func (t *CmapField[K, V]) Value() V {
	return t.v
}
