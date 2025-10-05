package once

import "sync"

type Map interface {
	Get(string) Once
	BoolOnce(string) bool
	ResetOnce(string)
	RunOnce(string, func())
	GoOnce(string, func())
}

var _ Map = (*onceMapImpl)(nil)

type onceMapImpl struct {
	mu    sync.RWMutex
	onces map[string]Once
}

func (o *onceMapImpl) Get(name string) Once {
	o.mu.RLock()
	defer o.mu.RUnlock()
	once, ok := o.onces[name]
	if !ok {
		o.mu.RUnlock()
		o.mu.Lock()
		once = New()
		o.onces[name] = once
		o.mu.Unlock()
		o.mu.RLock()
	}
	return once
}

func (o *onceMapImpl) BoolOnce(name string) bool     { return o.Get(name).Bool() }
func (o *onceMapImpl) ResetOnce(name string)         { o.Get(name).Reset() }
func (o *onceMapImpl) RunOnce(name string, f func()) { o.Get(name).Run(f) }
func (o *onceMapImpl) GoOnce(name string, f func())  { o.Get(name).Go(f) }

var defaultOnceMap Map

func NewMap() Map {
	return &onceMapImpl{onces: make(map[string]Once)}
}

func DefaultMap() Map {
	if defaultOnceMap == nil {
		defaultOnceMap = NewMap()
	}
	return defaultOnceMap
}
