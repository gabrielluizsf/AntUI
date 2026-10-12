//go:build android

package power

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// readBattery reads the sticky battery broadcast into out.
func readBattery(e *jni.Env, out *Battery) error {
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
	return readExtras(e, status, out)
}
