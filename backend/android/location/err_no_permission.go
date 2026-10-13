//go:build android

package location

import "errors"

// ErrNoPermission is asking for a position without having been granted one
// of the two location permissions.
var ErrNoPermission = errors.New("location: this app has neither " +
	"ACCESS_FINE_LOCATION nor ACCESS_COARSE_LOCATION; ask for one with " +
	"antui/backend/android/permission before asking where the device is")
