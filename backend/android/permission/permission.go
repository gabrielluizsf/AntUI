//go:build android

package permission

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// The permissions an app built with this library is most likely to want,
// spelled the way the manifest and the request both take them.
//
// Only the dangerous ones are here — the ones the user is asked about. A
// normal permission like INTERNET is granted at install by being in the
// manifest, and there is nothing to request.
const (
	Camera             = "android.permission.CAMERA"
	Microphone         = "android.permission.RECORD_AUDIO"
	FineLocation       = "android.permission.ACCESS_FINE_LOCATION"
	CoarseLocation     = "android.permission.ACCESS_COARSE_LOCATION"
	BackgroundLocation = "android.permission.ACCESS_BACKGROUND_LOCATION"
	ReadStorage        = "android.permission.READ_EXTERNAL_STORAGE"
	WriteStorage       = "android.permission.WRITE_EXTERNAL_STORAGE"
	ReadImages         = "android.permission.READ_MEDIA_IMAGES"
	ReadVideo          = "android.permission.READ_MEDIA_VIDEO"
	ReadAudio          = "android.permission.READ_MEDIA_AUDIO"
	Notifications      = "android.permission.POST_NOTIFICATIONS"
	Contacts           = "android.permission.READ_CONTACTS"
	Calendar           = "android.permission.READ_CALENDAR"
	BodySensors        = "android.permission.BODY_SENSORS"
	Activity           = "android.permission.ACTIVITY_RECOGNITION"
	NearbyDevices      = "android.permission.BLUETOOTH_CONNECT"
)

// Result is the answer to a request.
type Result = app.PermissionResult

// Held reports whether a permission has already been granted.
func Held(name string) (bool, error) { return app.HasPermission(name) }

// Request asks the user and waits.
//
// **Not from the frame loop.** The dialog pauses the app, Android waits for
// the app to acknowledge the pause, and a loop that is waiting here is a
// loop that is not acknowledging. Run it on a goroutine of its own.
func Request(names ...string) (Result, error) { return app.RequestPermissions(names...) }

// Ensure asks only for what is not already held, and reports whether
// everything asked for is held afterwards. It is what most callers want:
// asking again for something already granted shows no dialog, but it does
// pause the app for a moment, and doing it every launch is visible.
//
// **Not from the frame loop**, for [Request]'s reason — and the trap here is
// worse than Request's, because Ensure only blocks when something is
// actually missing. That is the *first* launch and no other. An app that
// calls this from its frame loop works perfectly every time its author runs
// it, and hangs on the one launch that matters: the user's first.
//
// Run it on a goroutine and let the loop keep drawing:
//
//	go func() { granted <- ensure() }()
func Ensure(names ...string) (bool, error) {
	var missing []string
	for _, n := range names {
		held, err := Held(n)
		if err != nil {
			return false, err
		}
		if !held {
			missing = append(missing, n)
		}
	}
	if len(missing) == 0 {
		return true, nil
	}
	r, err := Request(missing...)
	if err != nil {
		return false, err
	}
	return r.All(), nil
}

// ShouldExplain reports whether the user has refused this permission once
// and the app should say why it wants it before asking again.
//
// It is a three-state answer squeezed into two, and the third state is the
// one that matters: **false means either "never asked" or "refused for
// good"**. The platform does not distinguish them, deliberately — an app
// that could tell the difference could pester. So the way to use it is: ask;
// if refused and this is true, explain and ask once more; if refused and
// this is false, the dialog will not appear again and the only way forward
// is the settings screen.
func ShouldExplain(name string) (bool, error) {
	a := app.Current()
	if a == nil {
		return false, app.ErrNoUIThread
	}
	if a.Config().SDK < 23 {
		return false, nil
	}
	var out bool
	err := jni.Do(func(e *jni.Env) error {
		js, err := e.String(name)
		if err != nil {
			return err
		}
		out, err = e.InvokeBool(app.Context(), "shouldShowRequestPermissionRationale",
			jni.Sig(jni.TBool, jni.TString), jni.Ref(js))
		return err
	})
	return out, err
}

// OpenSettings takes the user to this app's settings page, which is the only
// way back from a permission refused for good.
func OpenSettings() error {
	return jni.Do(func(e *jni.Env) error {
		action, err := e.ConstantString("android/provider/Settings",
			"ACTION_APPLICATION_DETAILS_SETTINGS")
		if err != nil {
			return err
		}
		pkg, err := e.InvokeString(app.Context(), "getPackageName", jni.Sig(jni.TString))
		if err != nil {
			return err
		}
		jaction, err := e.String(action)
		if err != nil {
			return err
		}
		juri, err := e.String("package:" + pkg)
		if err != nil {
			return err
		}
		uri, err := e.Static("android/net/Uri", "parse",
			jni.Sig(jni.TClass("android/net/Uri"), jni.TString), jni.Ref(juri))
		if err != nil {
			return err
		}
		intent, err := e.Make("android/content/Intent",
			jni.Sig(jni.TVoid, jni.TString, jni.TClass("android/net/Uri")),
			jni.Ref(jaction), jni.Ref(uri))
		if err != nil {
			return err
		}
		return e.InvokeVoid(app.Context(), "startActivity",
			jni.Sig(jni.TVoid, jni.TClass("android/content/Intent")), jni.Ref(intent))
	})
}
