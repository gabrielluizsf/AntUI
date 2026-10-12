package apk

import (
	"fmt"
	"os"
	"path/filepath"
)

// checkRequired reports what the caller has to give: a package name and a
// Go package directory that exists.
func (c *Config) checkRequired() error {
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
	return nil
}
