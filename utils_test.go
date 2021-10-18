package watchdog

import "testing"

func TestCheckPath(t *testing.T) {
	err := checkPath("/var/temp/pprof")
	if err != nil {
		t.Error(err)
	}
}
