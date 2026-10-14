package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

// AssembleBundleParams holds the parameters for assembling an app bundle.
type AssembleBundleParams struct {
	Proto string
	Libs  map[sdk.ABI]string
	Out   string
}
