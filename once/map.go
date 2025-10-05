package once

import "sync"

type OnceMap interface {
	Get(string) Once
	BoolOnce(string) bool
	ResetOnce(string)
	RunOnce(string, func())
	GoOnce(string, func())
}

var _ OnceMap = (*onceMapImpl)(nil)

type onceMapImpl struct {
	mu    sync.RWMutex
	onces map[string]Once
}

func (o *onceMapImpl) Get(name string) Once {
	o.mu.RLock()
	defer o.mu.RUnlock()
	once, ok := o.onces[name]
	if !ok {
		o.mu.Lock()
		once = NewOnce()
		o.onces[name] = once
		o.mu.Unlock()
	}
	return once
}

func (o *onceMapImpl) BoolOnce(name string) bool     { return o.Get(name).Bool() }
func (o *onceMapImpl) ResetOnce(name string)         { o.Get(name).Reset() }
func (o *onceMapImpl) RunOnce(name string, f func()) { o.Get(name).Run(f) }
func (o *onceMapImpl) GoOnce(name string, f func())  { o.Get(name).Go(f) }
