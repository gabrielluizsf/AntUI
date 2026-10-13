//go:build android

package app

import "errors"

// ErrNoShim is what the calls here report when the app was built without the
// Java shim, so there is nothing to deliver the answer.
var ErrNoShim = errors.New("app: this app was built without the Java shim, so " +
	"permission and activity results cannot be delivered; build without -no-shim")
