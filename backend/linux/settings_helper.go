//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

import (
	"os/exec"
	"sync"
)

// settingsHelper is the program a desktop's settings are read from, looked
// up once: a machine without one is not asked about it again.
var settingsHelper = sync.OnceValue(func() string {
	path, err := exec.LookPath("gsettings")
	if err != nil {
		return ""
	}
	return path
})
