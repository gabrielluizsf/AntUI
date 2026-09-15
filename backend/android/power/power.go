//go:build android

package power

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Battery is what the device says about its own power.
type Battery struct {
	// Level is 0 to 1, or -1 when the device will not say.
	Level float32
	// Charging is whether it is going up.
	Charging bool
	// Full is charging that has finished.
	Full bool
	// Plugged is what it is plugged into: [Unplugged], [AC], [USB] or
	// [Wireless].
	Plugged int
}

// What the battery is plugged into.
const (
	Unplugged = 0
	AC        = 1
	USB       = 2
	Wireless  = 4
)

// The values BatteryManager gives for the charging status.
const (
	statusCharging = 2
	statusFull     = 5
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
		action, err := e.ConstantString("android/content/Intent", "ACTION_BATTERY_CHANGED")
		if err != nil {
			return err
		}
		jaction, err := e.String(action)
		if err != nil {
			return err
		}
		filter, err := e.Make("android/content/IntentFilter",
			jni.Sig(jni.TVoid, jni.TString), jni.Ref(jaction))
		if err != nil {
			return err
		}
		// A null receiver means "do not subscribe, just give me the last one
		// that was sent" — which for a sticky broadcast is the current state.
		status, err := e.Invoke(app.Context(), "registerReceiver",
			jni.Sig(jni.TClass("android/content/Intent"),
				jni.TClass("android/content/BroadcastReceiver"),
				jni.TClass("android/content/IntentFilter")),
			jni.Ref(jni.Object{}), jni.Ref(filter))
		if err != nil || status.IsNil() {
			return err
		}

		extra := func(name string) (int32, error) {
			k, err := e.String(name)
			if err != nil {
				return 0, err
			}
			return e.InvokeInt(status, "getIntExtra",
				jni.Sig(jni.TInt, jni.TString, jni.TInt), jni.Ref(k), jni.Int(-1))
		}
		level, err := extra("level")
		if err != nil {
			return err
		}
		scale, err := extra("scale")
		if err != nil {
			return err
		}
		if level >= 0 && scale > 0 {
			out.Level = float32(level) / float32(scale)
		}
		state, err := extra("status")
		if err != nil {
			return err
		}
		out.Charging = state == statusCharging
		out.Full = state == statusFull
		plugged, err := extra("plugged")
		if err != nil {
			return err
		}
		if plugged > 0 {
			out.Plugged = int(plugged)
		}
		return nil
	})
	return out, err
}

// Saving reports whether the user has turned battery saver on. An app that
// respects it does less: fewer frames, no background work, no animation
// that is only decoration.
func Saving() (bool, error) {
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServicePower)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isPowerSaveMode", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}

// Idle reports whether the device is in doze — the deep sleep it enters when
// it has been still and unplugged for a while, in which network access and
// alarms are held until it wakes.
//
// It is false below API 23, where there was no doze.
func Idle() (bool, error) {
	a := app.Current()
	if a == nil {
		return false, app.ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return false, nil
	}
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		manager, err := app.Service(e, app.ServicePower)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(manager, "isDeviceIdleMode", jni.Sig(jni.TBool))
		return err
	})
	return out, err
}
