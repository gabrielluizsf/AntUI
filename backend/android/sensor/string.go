//go:build android

package sensor

import (
	"fmt"
)

// String names the sensor, and says the number for one this package has
// no name for.
func (k Kind) String() string {
	if n, ok := kindNames[k]; ok {
		return n
	}
	return fmt.Sprintf("sensor(%d)", int(k))
}
