package hookutils

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type sigLoop struct {
	hooks map[syscall.Signal]*hook
	// Hook for all exit signals (SIGINT, SIGTERM, and SIGQUIT).
	// It is recommended to use this hook to register exit handlers
	// instead of adding a hook for each individual exit signal.
	exitHook *hook
	exiting  bool
	mu       sync.Mutex
}

func (s *sigLoop) Handle(sig syscall.Signal, h *hook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hooks[sig]; ok {
		return fmt.Errorf("hook for signal %d is already exists", sig)
	}
	s.hooks[sig] = h
	return nil
}

func (s *sigLoop) HandleExit(h *hook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.exitHook != nil {
		return errors.New("exit hook is already exists")
	}
	s.exitHook = h
	return nil
}

func (s *sigLoop) Get(sig syscall.Signal) *hook {
	if hook, ok := s.hooks[sig]; ok {
		return hook
	}
	return nil
}

func (s *sigLoop) GetExit() *hook {
	return s.exitHook
}

func (s *sigLoop) Loop() {
	sigchan := make(chan os.Signal, 1)
	for sig := range s.hooks {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		default:
			signal.Notify(sigchan, sig)
		}
	}
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		sig := <-sigchan
		if s.exiting {
			return
		}
		if hook, ok := s.hooks[sig.(syscall.Signal)]; ok {
			hook.Exec()
		}
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			s.Exit(0)
		}
	}
}

func (s *sigLoop) Exit(code int) {
	if s.exiting {
		return
	}
	s.mu.Lock()
	s.exiting = true
	s.mu.Unlock()
	if s.exitHook != nil {
		s.exitHook.Exec()
	}
	for sig := range s.hooks {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
		default:
			signal.Reset(sig)
		}
	}
	signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	os.Exit(code)
}

var defaultSigLoop *sigLoop

func NewSigLoop() *sigLoop {
	return &sigLoop{hooks: make(map[syscall.Signal]*hook)}
}

func DefaultSigLoop() *sigLoop {
	if defaultSigLoop == nil {
		defaultSigLoop = NewSigLoop()
	}
	return defaultSigLoop
}
