package apk

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrielluizsf/antui/canvas"
)

// checkOptional fills the defaults for the optional parts and checks that
// the icon and the assets directory are where the caller said.
func (c *Config) checkOptional() error {
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
	return nil
}
