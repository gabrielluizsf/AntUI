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
	// What the store would say, said now. A bundle is built to be uploaded,
	// so it is worth hearing before the two minutes of building rather than
	// after — but it is said and not enforced: a bundle built to be tried on
	// an emulator wants x86_64, which no phone has and no store accepts, and
	// that is a perfectly good reason to build one.
	for _, p := range cfg.PlayProblems() {
		fmt.Fprintln(os.Stderr, "antuiapk: "+p.String())
	}

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

	st, err := os.Stat(cfg.Out)
	if err != nil {
		return nil, fmt.Errorf("apk: %w", err)
	}
	res := &Result{APK: cfg.Out, Libs: libs, Size: st.Size()}
	// The bundle already carries the symbols inside it, so a copy beside it
	// would be the same three megabytes twice. One is written only when the
	// caller asked for one by name.
	if cfg.Symbols != "" {
		if wrote, err := WriteSymbols(libs, cfg.Lib, cfg.Symbols); err != nil {
			return nil, err
		} else if wrote {
			res.Symbols = cfg.Symbols
		}
	}
	return res, nil
}
