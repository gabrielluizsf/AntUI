package sdk

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// readNDK reads an NDK's version out of source.properties and finds the
// prebuilt toolchain built for this machine.
func readNDK(dir string) (NDK, bool) {
	b, err := os.ReadFile(filepath.Join(dir, "source.properties"))
	if err != nil {
		return NDK{}, false
	}
	n := NDK{Dir: dir}
	for line := range strings.Lines(string(b)) {
		if rev, ok := strings.CutPrefix(strings.TrimSpace(line), "Pkg.Revision"); ok {
			if _, v, found := strings.Cut(rev, "="); found {
				n.Version = ParseVersion(strings.TrimSpace(v))
			}
		}
	}
	prebuilt := filepath.Join(dir, "toolchains", "llvm", "prebuilt")
	hosts := subdirs(prebuilt)
	if len(hosts) == 0 {
		return NDK{}, false
	}
	// Take the directory whose name starts with this OS. The NDK has shipped
	// darwin-x86_64 on Apple Silicon and darwin-arm64 in different revisions,
	// so the exact name is read rather than assumed, and any host directory
	// is better than refusing to build.
	n.Host = hosts[0]
	for _, h := range hosts {
		if strings.HasPrefix(h, runtime.GOOS) {
			n.Host = h
			break
		}
	}
	return n, n.Version.Major() > 0
}
