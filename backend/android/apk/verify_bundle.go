package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// VerifyBundle takes a bundle apart the way the store will, and puts the
// result on a device.
//
// This is the only check that means anything about a bundle. Everything
// before it says the file is shaped right; this says the store can turn it
// into an app that installs and runs — which is a different question, and
// the one that matters at three in the morning before a release.
//
// serial names a device, or is empty for the only one attached.
func VerifyBundle(tc *sdk.Toolchain, cfg *Config, params VerifyBundleParams) error {
	bundle := params.Bundle
	serial := params.Serial
	if err := cfg.check(); err != nil {
		return err
	}
	if tc.Bundletool == "" {
		return ErrNoBundletool
	}
	if tc.JDK.Dir == "" {
		return fmt.Errorf("apk: bundletool is a jar and needs a JDK to run")
	}

	work, err := os.MkdirTemp("", "antuiverify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	apks := filepath.Join(work, "out.apks")

	if err := buildApksFromBundle(tc, cfg, bundle, serial, apks); err != nil {
		return err
	}
	return installApks(tc, serial, apks)
}
