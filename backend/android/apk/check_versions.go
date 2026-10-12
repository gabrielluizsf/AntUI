package apk

import (
	"fmt"
	"path/filepath"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// checkVersions fills the labels and the version and SDK numbers, and
// refuses a minSdk above the target.
func (c *Config) checkVersions() error {
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
	return nil
}
