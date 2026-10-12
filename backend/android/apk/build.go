package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Build compiles the Go package and produces a signed, installable APK.
//
// It is the four steps below run in order, in a temporary directory that is
// removed afterwards. Nothing is written outside it except the APK itself.
func Build(tc *sdk.Toolchain, cfg Config) (*Result, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	if tc.NDK.Dir == "" {
		return nil, fmt.Errorf("apk: no NDK; a Go app is native code and cannot be built without one")
	}
	if tc.BuildTools.Dir == "" || tc.Platform.Dir == "" {
		return nil, fmt.Errorf("apk: the SDK has no build-tools or no platform installed")
	}

	work, err := os.MkdirTemp("", "antuiapk-")
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	defer os.RemoveAll(work)

	libs, err := BuildLibs(tc, &cfg, work)
	if err != nil {
		return nil, err
	}
	base := filepath.Join(work, "base.apk")
	if err := Link(tc, &cfg, base); err != nil {
		return nil, err
	}
	packed := filepath.Join(work, "packed.apk")
	if err := Assemble(base, libs, cfg.Lib, packed, &cfg); err != nil {
		return nil, err
	}
	if err := Sign(tc, &cfg, packed, cfg.Out); err != nil {
		return nil, err
	}
	return buildResult(cfg, libs)
}
