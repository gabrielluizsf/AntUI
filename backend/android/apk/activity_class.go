package apk

import (
	"github.com/gabrielluizsf/antui/backend/android/shim"
)

// ActivityClass is the activity the manifest names: the shim when there is
// one, and the platform's own NativeActivity when there is not.
func (c *Config) ActivityClass() string {
	// nil is the default and the default is on, which is what [Config.check]
	// fills in. Answering "NativeActivity" for nil here made this disagree
	// with the manifest that check produces, so a caller who asked before
	// building — anything that installs and then starts the app — started an
	// activity the package does not declare.
	if c.Shim == nil || *c.Shim {
		return shim.Activity
	}
	return "android.app.NativeActivity"
}
