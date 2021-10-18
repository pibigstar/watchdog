package watchdog_test

import (
	"github.com/pibigstar/watchdog"
	"testing"
)

func TestGoroutineWatcher(t *testing.T) {
	w := watchdog.NewGoroutineWatcher()
	watchdog.Run(w)
}

func TestNewCpuWatcher(t *testing.T) {
	w := watchdog.NewCpuWatcher()
	watchdog.Run(w)
}

func TestNewMemoryWatcher(t *testing.T) {
	w := watchdog.NewMemoryWatcher()
	watchdog.Run(w)
}
