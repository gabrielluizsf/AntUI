package shim

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Build compiles the shim from [Source] and returns the dex.
//
// It needs a JDK and the SDK's build tools, which building an app does not —
// which is the whole reason the result is committed. Nothing is written
// outside the temporary directory it makes.
func Build(tc *sdk.Toolchain) ([]byte, error) {
	if tc.JDK.Javac == "" {
		return nil, fmt.Errorf("shim: rebuilding needs a JDK, and none was found")
	}
	if tc.Platform.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs an SDK platform to compile against")
	}
	if tc.BuildTools.Dir == "" {
		return nil, fmt.Errorf("shim: rebuilding needs the SDK build tools for d8")
	}

	work, err := os.MkdirTemp("", "antui-shim-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)

	srcs, err := writeSources(work)
	if err != nil {
		return nil, err
	}
	classes := filepath.Join(work, "classes")
	if err := os.MkdirAll(classes, 0o755); err != nil {
		return nil, err
	}
	if err := compileJava(tc, classes, srcs); err != nil {
		return nil, err
	}
	compiled, err := classFiles(classes)
	if err != nil {
		return nil, err
	}
	if err := runD8(tc, work, compiled); err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(work, "classes.dex"))
}
