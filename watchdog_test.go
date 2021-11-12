package watchdog_test

import (
	"github.com/pibigstar/watchdog"
	"testing"
	"time"
)

func TestGoroutineWatcher(t *testing.T) {
	w := watchdog.NewGoroutineWatcher()
	w.MaxThreshold = 3
	watchdog.Run(w)

	for i:=0;i<10;i++ {
		go func() {
			time.Sleep(10 * time.Second)
		}()

		time.Sleep(1 * time.Second)
	}

	select {}
}

func TestNewCpuWatcher(t *testing.T) {
	w := watchdog.NewCpuWatcher()
	watchdog.Run(w)
}

func TestNewMemoryWatcher(t *testing.T) {
	w := watchdog.NewMemoryWatcher()
	watchdog.Run(w)
}
