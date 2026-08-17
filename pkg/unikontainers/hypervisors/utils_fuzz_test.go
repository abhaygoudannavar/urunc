package hypervisors

import (
	"testing"
)

func FuzzBytesToStringMB(f *testing.F) {
	f.Add(uint64(256000000))
	f.Add(uint64(0))
	f.Add(uint64(1))
	f.Add(uint64(1048576))

	f.Fuzz(func(t *testing.T, argMem uint64) {
		res := BytesToStringMB(argMem)
		if res == "" {
			t.Errorf("BytesToStringMB returned empty string for %d", argMem)
		}
	})
}
