package once

import "sync"

type Once interface {
	Bool() bool
	Reset()
	Run(func())
	Go(func())
}

var _ Once = (*onceImpl)(nil)

type onceImpl struct {
	mu   sync.Mutex
	cond bool
}

func (o *onceImpl) Bool() bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.cond {
		o.cond = false
		return true
	}
	return false
}

func (o *onceImpl) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.cond = true
}

func (o *onceImpl) Run(f func()) {
	if o.Bool() {
		f()
	}
}

func (o *onceImpl) Go(f func()) {
	if o.Bool() {
		go f()
	}
}

func New() Once {
	return &onceImpl{cond: true}
}
