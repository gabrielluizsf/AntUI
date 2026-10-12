package apk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// goCommand finds the go tool: on PATH first, and in GOROOT when the build
// is running somewhere PATH was never set up, which is most CI.
func goCommand() (string, error) {
	if p, err := exec.LookPath("go"); err == nil {
		return p, nil
	}
	p := filepath.Join(runtime.GOROOT(), "bin", "go")
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}
	return "", fmt.Errorf("apk: no go tool on PATH")
}
