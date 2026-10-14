package apk

import "github.com/gabrielluizsf/antui/backend/android/sdk"

// AssembleParams holds the parameters for assembling an APK.
type AssembleParams struct {
	Base    string
	Libs    map[sdk.ABI]string
	LibName string
	Out     string
}
