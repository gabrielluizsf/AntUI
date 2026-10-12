package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// Bundle builds an .aab, signed and ready to upload.
func Bundle(tc *sdk.Toolchain, cfg Config) (*Result, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	if tc.NDK.Dir == "" {
		return nil, fmt.Errorf("apk: no NDK; a Go app is native code")
	}
	if tc.JDK.Dir == "" {
		return nil, fmt.Errorf("apk: signing a bundle needs a JDK: apksigner " +
			"cannot sign one, and jarsigner is what does")
	}
	if cfg.Out == cfg.Label+".apk" {
		cfg.Out = cfg.Label + ".aab"
	}
	warnPlayProblems(cfg)

	work, err := os.MkdirTemp("", "antuiaab-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(work)

	libs, err := BuildLibs(tc, &cfg, work)
	if err != nil {
		return nil, err
	}
	proto := filepath.Join(work, "proto.apk")
	if err := LinkProto(tc, &cfg, proto); err != nil {
		return nil, err
	}
	unsigned := filepath.Join(work, "unsigned.aab")
	if err := AssembleBundle(proto, libs, &cfg, unsigned); err != nil {
		return nil, err
	}
	if err := SignBundle(tc, &cfg, unsigned, cfg.Out); err != nil {
		return nil, err
	}
	return bundleResult(cfg, libs)
}
