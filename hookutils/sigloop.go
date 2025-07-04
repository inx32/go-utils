package hookutils

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type sigLoop struct {
	hooks map[syscall.Signal]*hook
	mu    sync.Mutex
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

func (s *sigLoop) Get(sig syscall.Signal) *hook {
	if hook, ok := s.hooks[sig]; ok {
		return hook
	}
	return nil
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
		if hook, ok := s.hooks[sig.(syscall.Signal)]; ok {
			hook.Exec()
		}
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			for sig := range s.hooks {
				switch sig {
				case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				default:
					signal.Reset(sig)
				}
			}
			signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			os.Exit(0)
		}
	}
}

func NewSigLoop() *sigLoop {
	return &sigLoop{hooks: make(map[syscall.Signal]*hook)}
}
