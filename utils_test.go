package watchdog

import "testing"

func TestCheckPath(t *testing.T) {
	err := checkPath("/var/tmp/pprof")
	if err != nil {
		t.Error(err)
	}
}

func TestGenPassword(t *testing.T) {
	s := "test"
	t.Log(GenPassword(s))
}
