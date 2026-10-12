//go:build android

package assets

import (
	"github.com/gabrielluizsf/antui/backend/android/ndk"
)

// file is one asset open for reading.
type file struct {
	name  string
	asset *ndk.Asset
}
