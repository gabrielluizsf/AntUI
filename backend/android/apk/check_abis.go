package apk

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/sdk"
)

// checkABIs fills the ABI list when it is empty and refuses one the platform
// does not sell.
func (c *Config) checkABIs() error {
	if len(c.ABIs) == 0 {
		c.ABIs = sdk.Ship
	}
	for _, a := range c.ABIs {
		if !a.Valid() {
			return fmt.Errorf("apk: %q is not an Android ABI", a)
		}
	}
	return nil
}
