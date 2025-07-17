package cmap

import "sync"

// Thread-safe map
type Cmap[K comparable, V any] interface {
	Set(key K, value V)
	SetFunc(key K, f func(value V, ex bool) V)
	Remove(key K)
	Get(key K) (value V, ex bool)
	GetDefault(key K, def V) (value V)
	Has(key K) (ex bool)
	Keys() <-chan K
	Values() <-chan V
	Range() <-chan CmapField[K, V]
}

type cmapImpl[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func (t *cmapImpl[K, V]) Set(k K, v V) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.m[k] = v
}

func (t *cmapImpl[K, V]) SetFunc(k K, f func(V, bool) V) {
	val, ok := t.Get(k)
	t.Set(k, f(val, ok))
}

func (t *cmapImpl[K, V]) Remove(k K) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.m, k)
}

func (t *cmapImpl[K, V]) Get(k K) (v V, ok bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if val, ok := t.m[k]; ok {
		return val, true
	}
	return v, false
}

func (t *cmapImpl[K, V]) GetDefault(k K, def V) V {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if val, ok := t.m[k]; ok {
		return val
	}
	return def
}

func (t *cmapImpl[K, V]) Has(k K) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.m[k]
	return ok
}

func (t *cmapImpl[K, V]) Keys() <-chan K {
	c := make(chan K)
	go func() {
		t.mu.RLock()
		defer t.mu.RUnlock()
		for k := range t.m {
			c <- k
		}
		close(c)
	}()
	return c
}

func (t *cmapImpl[K, V]) Values() <-chan V {
	c := make(chan V)
	go func() {
		t.mu.RLock()
		defer t.mu.RUnlock()
		for _, v := range t.m {
			c <- v
		}
		close(c)
	}()
	return c
}

func (t *cmapImpl[K, V]) Range() <-chan CmapField[K, V] {
	c := make(chan CmapField[K, V])
	go func() {
		t.mu.RLock()
		defer t.mu.RUnlock()
		for k, v := range t.m {
			c <- &cmapFieldImpl[K, V]{k, v}
		}
		close(c)
	}()
	return c
}

func New[K comparable, V any]() Cmap[K, V] {
	return &cmapImpl[K, V]{m: make(map[K]V)}
}
