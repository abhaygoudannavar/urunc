package unikontainers

import (
	"testing"
)

func FuzzUruncConfigFromMap(f *testing.F) {
	f.Add("urunc_config.monitors.qemu.default_memory_mb", "256")
	f.Add("urunc_config.extra_binaries.virtiofsd.path", "/usr/libexec/virtiofsd")

	f.Fuzz(func(t *testing.T, fuzzKey string, fuzzVal string) {
		cfgMap := map[string]string{
			fuzzKey: fuzzVal,
		}
		_ = UruncConfigFromMap(cfgMap)
	})
}
