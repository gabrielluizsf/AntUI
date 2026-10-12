//go:build android

package camera

import (
	"fmt"
)

// String describes the running camera: which one, and how its frames are
// turned to be the right way up.
func (p *Preview) String() string {
	return fmt.Sprintf("%s camera %s, turned %d°", p.info.Facing, p.info.ID, p.upright)
}
