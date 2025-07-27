package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SigLoop interface {
	// Handle registers a [Hook] for specific signal.
	// If a hook for a signal has already been registered, Handle will return an error.
	Handle(sig syscall.Signal, h Hook) error

	// HandleExit registers a program exit [Hook].
	// It is better to register a program exit hook than to register hook for each signal.
	// If the exit hook has already been registered, HandleExit will return an error.
	HandleExit(h Hook) error

	// Get returns [Hook] for the signal.
	// If the hook for the signal is not registered, Get will return <nil>.
	Get(sig syscall.Signal) Hook

	// GetOrHandle returns a [Hook] for the signal.
	// Unlike Get, GetOrHandle registers a hook for a signal if not registered.
	GetOrHandle(sig syscall.Signal) Hook

	// Get returns a program exit [Hook].
	// If the exit hook is not registered, GetExit will return <nil>.
	GetExit() Hook

	// GetExitOrHandle returns a program exit [Hook].
	// Unlike GetExit, GetExitOrHandle registers a exit hook if not registered.
	GetExitOrHandle() Hook

	// Loop starts listening for signals or Exit.
	// If an exit signal is received (SIGINT, SIGTERM or SIGQUIT), the program will exit.
	Loop()

	// Exit runs the exit [Hook] and stops the program with a specified code.
	// Exit will NOT trigger any signals. Use HandleExit to handle program exit.
	Exit(code int)
}

var _ SigLoop = (*sigLoopImpl)(nil)

type sigLoopImpl struct {
	hooks map[syscall.Signal]Hook
	// Hook for all exit signals (SIGINT, SIGTERM, and SIGQUIT).
	// It is recommended to use this hook to register exit handlers
	// instead of adding a hook for each individual exit signal.
	exitHook Hook
	exiting  bool
	mu       sync.Mutex
}

func (s *sigLoopImpl) Handle(sig syscall.Signal, h Hook) error {
	if h == nil {
		return errors.New("hook is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hooks[sig]; ok {
		return fmt.Errorf("hook for signal %d is already exists", sig)
	}
	s.hooks[sig] = h
	return nil
}

func (s *sigLoopImpl) HandleExit(h Hook) error {
	if h == nil {
		return errors.New("hook is nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.exitHook != nil {
		return errors.New("exit hook is already exists")
	}
	s.exitHook = h
	return nil
}

func (s *sigLoopImpl) Get(sig syscall.Signal) Hook {
	if hook, ok := s.hooks[sig]; ok {
		return hook
	}
	return nil
}

func (s *sigLoopImpl) GetOrHandle(sig syscall.Signal) Hook {
	h := s.Get(sig)
	if h == nil {
		h = NewHook(sig.String(), fmt.Sprintf("Handle signal %d \"%s\"", sig, sig.String()))
		s.Handle(sig, h)
	}
	return h
}

func (s *sigLoopImpl) GetExit() Hook {
	return s.exitHook
}

func (s *sigLoopImpl) GetExitOrHandle() Hook {
	if s.exitHook != nil {
		return s.exitHook
	}
	s.HandleExit(NewHook("exit", "Handle exit"))
	return s.exitHook
}

func (s *sigLoopImpl) Loop() {
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

func (s *sigLoopImpl) Exit(code int) {
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

var defaultSigLoop SigLoop

func NewSigLoop() SigLoop {
	return &sigLoopImpl{hooks: make(map[syscall.Signal]Hook)}
}

func DefaultSigLoop() SigLoop {
	if defaultSigLoop == nil {
		defaultSigLoop = NewSigLoop()
	}
	return defaultSigLoop
}
