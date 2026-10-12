//go:build android

package power

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Read asks the battery how it is.
//
// It goes through the sticky broadcast the system keeps rather than through
// BatteryManager's own getters, because that broadcast is API 1 and reports
// everything at once, while the getters arrived at three different API
// levels and each is one call.
func Read() (Battery, error) {
	out := Battery{Level: -1}
	err := jni.Do(func(e *jni.Env) error {
		return readBattery(e, &out)
	})
	return out, err
}
