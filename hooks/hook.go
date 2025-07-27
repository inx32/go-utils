package hooks

import (
	"errors"
	"fmt"
	"slices"
	"sync"
)

type Hook interface {
	// Exec runs the registered notifies and functions in order of decreasing weight.
	// First, Exec runs notifies, then it runs functions.
	Exec()

	// Handle registers a function in a hook.
	// If the function name already registered, Handle returns an error.
	Handle(*HookFunc) error

	// Notify registers a notify in a hook.
	// If the notify name already registered, Notify returns an error.
	Notify(*HookNotify) error
}

var _ Hook = (*hookImpl)(nil)

type hookImpl struct {
	name        string
	desc        string
	funcList    []*HookFunc
	notifyList  []*HookNotify
	funcNames   map[string]struct{}
	notifyNames map[string]struct{}
	mu          sync.Mutex
}

func (h *hookImpl) ensureSorted() {
	var last uint16
	for _, v := range h.funcList {
		if last < v.Weight {
			h.sort()
			return
		}
		last = v.Weight
	}
	last = 0
	for _, v := range h.notifyList {
		if last < v.Weight {
			h.sort()
			return
		}
		last = v.Weight
	}
}

func (h *hookImpl) sort() {
	h.mu.Lock()
	defer h.mu.Unlock()
	slices.SortFunc(h.funcList, func(a, b *HookFunc) int {
		return int(b.Weight) - int(a.Weight)
	})
	slices.SortFunc(h.notifyList, func(a, b *HookNotify) int {
		return int(b.Weight) - int(a.Weight)
	})
}

func (h *hookImpl) Exec() {
	h.ensureSorted()

	for _, n := range h.notifyList {
		if n.NonBlocking {
			select {
			case n.Chan <- struct{}{}:
			default:
			}
		} else {
			n.Chan <- struct{}{}
		}

		if n.DoneChan != nil {
			<-n.DoneChan
		}
	}

	for _, f := range h.funcList {
		if f.Concurrent {
			go f.Func()
		} else {
			f.Func()
		}
	}
}

func (h *hookImpl) Handle(f *HookFunc) error {
	if f == nil {
		return errors.New("HookFunc is nil")
	}
	if f.Func == nil {
		return errors.New("func is nil")
	}
	if f.Name == "" {
		return errors.New("name is empty")
	}

	if _, ok := h.funcNames[f.Name]; ok {
		return fmt.Errorf("handler named \"%s\" is already exists", f.Name)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.funcNames[f.Name] = struct{}{}
	h.funcList = append(h.funcList, f)
	return nil
}

func (h *hookImpl) Notify(n *HookNotify) error {
	if n == nil {
		return errors.New("HookNotify is nil")
	}
	if n.Chan == nil {
		return errors.New("chan is nil")
	}
	if n.Name == "" {
		return errors.New("name is empty")
	}

	if _, ok := h.notifyNames[n.Name]; ok {
		return fmt.Errorf("notify named \"%s\" is already exists", n.Name)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notifyNames[n.Name] = struct{}{}
	h.notifyList = append(h.notifyList, n)
	return nil
}

func NewHook(name, desc string) Hook {
	return &hookImpl{
		name: name, desc: desc,
		funcNames:   make(map[string]struct{}),
		notifyNames: make(map[string]struct{}),
	}
}
