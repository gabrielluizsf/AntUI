//go:build android

package ndk

import "errors"

// ErrNoAsset is a name that is not in the package.
var ErrNoAsset = errors.New("ndk: no such asset")
