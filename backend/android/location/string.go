//go:build android

package location

import (
	"fmt"
)

// String is the position and how good it is, in one line.
func (f Fix) String() string {
	return fmt.Sprintf("%.6f, %.6f ±%.0fm (%s)", f.Latitude, f.Longitude, f.Accuracy, f.Provider)
}
