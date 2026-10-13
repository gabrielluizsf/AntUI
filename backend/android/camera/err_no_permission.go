//go:build android

package camera

import "errors"

// ErrNoPermission is opening a camera without having been granted one.
var ErrNoPermission = errors.New("camera: this app has not been granted CAMERA; " +
	"ask for it with antui/backend/android/permission first")
