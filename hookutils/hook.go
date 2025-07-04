package hookutils

import (
	"fmt"
	"slices"
	"sync"
)

type hook struct {
	name       string
	desc       string
	funcList   []*HookFunc
	notifyList []*HookNotify
	funcReg    map[string]struct{}
	notifyReg  map[string]struct{}
	mu         sync.Mutex
}

func (h *hook) ensureSorted() {
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

func (h *hook) sort() {
	h.mu.Lock()
	defer h.mu.Unlock()
	slices.SortFunc(h.funcList, func(a, b *HookFunc) int {
		return int(b.Weight) - int(a.Weight)
	})
	slices.SortFunc(h.notifyList, func(a, b *HookNotify) int {
		return int(b.Weight) - int(a.Weight)
	})
}

func (h *hook) Exec() {
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
			// TODO: handle error
			f.Func()
		}
	}
}

func (h *hook) Handle(f *HookFunc) error {
	if _, ok := h.funcReg[f.Name]; ok {
		return fmt.Errorf("handler named \"%s\" is already exists", f.Name)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.funcReg[f.Name] = struct{}{}
	h.funcList = append(h.funcList, f)
	return nil
}

func (h *hook) Notify(n *HookNotify) error {
	if _, ok := h.notifyReg[n.Name]; ok {
		return fmt.Errorf("notify named \"%s\" is already exists", n.Name)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notifyReg[n.Name] = struct{}{}
	h.notifyList = append(h.notifyList, n)
	return nil
}

func NewHook(name, desc string) *hook {
	return &hook{
		name: name, desc: desc,
		funcReg:   make(map[string]struct{}),
		notifyReg: make(map[string]struct{}),
	}
}
