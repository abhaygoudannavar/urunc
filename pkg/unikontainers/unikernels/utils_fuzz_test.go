package unikernels

import (
	"testing"
)

func FuzzSubnetMaskToCIDR(f *testing.F) {
	f.Add("255.255.255.0")
	f.Add("255.0.0.0")
	f.Add("0.0.0.0")

	f.Fuzz(func(t *testing.T, fuzzMask string) {
		_, _ = subnetMaskToCIDR(fuzzMask)
	})
}
