package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
	"github.com/gabrielluizsf/antui/canvas"
)

// check fills in the defaults and reports what is still missing. It is
// called by every entry point, so a caller who skips Build still gets it.
func (c *Config) check() error {
	if c.Package == "" {
		return fmt.Errorf("apk: no package name; it looks like \"com.example.app\"")
	}
	if err := validPackage(c.Package); err != nil {
		return err
	}
	if c.Dir == "" {
		return fmt.Errorf("apk: no Go package to build")
	}
	abs, err := filepath.Abs(c.Dir)
	if err != nil {
		return fmt.Errorf("apk: %w", err)
	}
	c.Dir = abs
	if st, err := os.Stat(c.Dir); err != nil || !st.IsDir() {
		return fmt.Errorf("apk: %s is not a directory", c.Dir)
	}

	if c.Label == "" {
		c.Label = filepath.Base(c.Dir)
	}
	if c.Lib == "" {
		c.Lib = defaultLib
	}
	if c.VersionCode == 0 {
		c.VersionCode = 1
	}
	if c.VersionName == "" {
		c.VersionName = defaultVersion
	}
	if c.MinSDK == 0 {
		c.MinSDK = sdk.MinSDK
	}
	if c.TargetSDK == 0 {
		c.TargetSDK = sdk.PlayTarget
	}
	if c.MinSDK > c.TargetSDK {
		return fmt.Errorf("apk: minSdk %d is above targetSdk %d", c.MinSDK, c.TargetSDK)
	}
	if len(c.ABIs) == 0 {
		c.ABIs = sdk.Ship
	}
	for _, a := range c.ABIs {
		if !a.Valid() {
			return fmt.Errorf("apk: %q is not an Android ABI", a)
		}
	}
	for _, l := range c.Links {
		if err := l.valid(); err != nil {
			return err
		}
	}
	if c.Strip == nil {
		on := true
		c.Strip = &on
	}
	if c.Shim == nil {
		on := true
		c.Shim = &on
	}
	if c.Store == nil {
		c.Store = DefaultStore
	}
	if c.Icon != "" {
		abs, err := filepath.Abs(c.Icon)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		if _, err := os.Stat(abs); err != nil {
			return fmt.Errorf("apk: there is no icon at %s", c.Icon)
		}
		c.Icon = abs
	}
	if c.IconBackground == 0 {
		c.IconBackground = canvas.White
	}
	if c.Assets != "" {
		abs, err := filepath.Abs(c.Assets)
		if err != nil {
			return fmt.Errorf("apk: %w", err)
		}
		if st, err := os.Stat(abs); err != nil || !st.IsDir() {
			return fmt.Errorf("apk: %s is not a directory of assets", c.Assets)
		}
		c.Assets = abs
	}
	if c.Out == "" {
		c.Out = c.Label + ".apk"
	}
	c.Keystore.fill()
	return nil
}
