//go:build android

package app

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Insets is how much of the window the system is covering, in pixels.
//
// A fullscreen app gets a surface the whole size of the screen and the
// system then draws its bars **over** it — so the app owns every pixel and
// some of them are not visible. This is which ones.
type Insets struct {
	Left, Top, Right, Bottom int
	// Keyboard is how tall the on-screen keyboard is, and 0 when it is
	// down. It is reported separately because it comes and goes while the
	// bars do not, and because a layout usually wants to move rather than
	// shrink for it.
	//
	// It is only available from API 30. Below that the keyboard is included
	// in Bottom instead, and there is no way to tell the two apart.
	Keyboard int
}

// WindowInsets asks the system what it is covering.
//
// It runs on the UI thread, and waits — a view may not be read from anywhere
// else. Call it when something changes rather than every frame; the events
// that change it are ContentRect and ConfigChanged.
func WindowInsets() (Insets, error) {
	var out Insets
	var err error
	if e := RunOnUISync(func() { out, err = readInsets() }); e != nil {
		return out, e
	}
	return out, err
}

func readInsets() (Insets, error) {
	var out Insets
	err := jni.Do(func(e *jni.Env) error {
		decor, err := decorView(e)
		if err != nil {
			return err
		}
		viewCls, err := e.Class("android/view/View")
		if err != nil {
			return err
		}
		// getRootWindowInsets arrived in API 23. Below that there is no way
		// to ask at all, so the honest answer is nothing rather than a guess.
		get, err := e.Method(viewCls, "getRootWindowInsets",
			jni.Sig(jni.TClass("android/view/WindowInsets")))
		if err != nil {
			return nil
		}
		ins, err := e.CallObject(decor, get)
		if err != nil || ins.IsNil() {
			return err
		}
		wiCls, err := e.Class("android/view/WindowInsets")
		if err != nil {
			return err
		}
		for _, side := range []struct {
			name string
			into *int
		}{
			{"getSystemWindowInsetLeft", &out.Left},
			{"getSystemWindowInsetTop", &out.Top},
			{"getSystemWindowInsetRight", &out.Right},
			{"getSystemWindowInsetBottom", &out.Bottom},
		} {
			m, err := e.Method(wiCls, side.name, jni.Sig(jni.TInt))
			if err != nil {
				return err
			}
			v, err := e.CallInt(ins, m)
			if err != nil {
				return err
			}
			*side.into = int(v)
		}
		out.Keyboard = imeInset(e, ins, wiCls)
		return nil
	})
	return out, err
}

// imeInset reads the keyboard's height on the platforms that report it
// separately. It answers 0 rather than an error on the ones that do not,
// because "there is no such call here" is not a failure.
func imeInset(e *jni.Env, ins jni.Object, wiCls jni.Class) int {
	typeCls, err := e.Class("android/view/WindowInsets$Type")
	if err != nil {
		return 0
	}
	imeM, err := e.StaticMethod(typeCls, "ime", jni.Sig(jni.TInt))
	if err != nil {
		return 0
	}
	which, err := e.CallStaticInt(typeCls, imeM)
	if err != nil {
		return 0
	}
	getInsets, err := e.Method(wiCls, "getInsets",
		jni.Sig(jni.TClass("android/graphics/Insets"), jni.TInt))
	if err != nil {
		return 0
	}
	got, err := e.CallObject(ins, getInsets, jni.Int(which))
	if err != nil || got.IsNil() {
		return 0
	}
	insetsCls, err := e.Class("android/graphics/Insets")
	if err != nil {
		return 0
	}
	f, err := e.Field(insetsCls, "bottom", jni.TInt)
	if err != nil {
		return 0
	}
	v, err := e.GetInt(got, f)
	if err != nil {
		return 0
	}
	return int(v)
}

// The system-UI flags, as android.view.View numbers them. They are
// deprecated in favour of WindowInsetsController from API 30 and still work,
// which is why they are what is used: one path that works everywhere beats
// two that each work half the time.
const (
	uiFullscreen       = 0x00000004 // hide the status bar
	uiHideNavigation   = 0x00000002 // hide the navigation bar
	uiImmersiveSticky  = 0x00001000 // let a swipe bring them back briefly
	uiLayoutStable     = 0x00000100
	uiLayoutHideNav    = 0x00000200
	uiLayoutFullscreen = 0x00000400
)

// SetImmersive hides the status and navigation bars, giving the app the
// whole screen. A swipe from an edge brings them back for a few seconds and
// then they go again, which is what "sticky" means and what a game wants.
func SetImmersive(on bool) error {
	flags := 0
	if on {
		flags = uiFullscreen | uiHideNavigation | uiImmersiveSticky |
			uiLayoutStable | uiLayoutHideNav | uiLayoutFullscreen
	}
	return setSystemUI(flags)
}

// SetEdgeToEdge draws the app behind the system bars without hiding them —
// the bars stay, over the app's own pixels, and [WindowInsets] is how the
// layout keeps out from under them.
func SetEdgeToEdge(on bool) error {
	flags := 0
	if on {
		flags = uiLayoutStable | uiLayoutHideNav | uiLayoutFullscreen
	}
	return setSystemUI(flags)
}

func setSystemUI(flags int) error {
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			decor, err := decorView(e)
			if err != nil {
				return err
			}
			viewCls, err := e.Class("android/view/View")
			if err != nil {
				return err
			}
			m, err := e.Method(viewCls, "setSystemUiVisibility", jni.Sig(jni.TVoid, jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(decor, m, jni.Int(int32(flags)))
		})
	}); e != nil {
		return e
	}
	return err
}

// The values android.content.pm.ActivityInfo uses for a requested
// orientation. Only the four worth asking for are here.
const (
	orientationUnspecified = -1
	orientationLandscape   = 0
	orientationPortrait    = 1
	orientationLocked      = 14 // whichever way it is now, and stay there
)

// SetOrientation asks the system to turn the screen and keep it there.
// [OrientationUnknown] gives the choice back to the device.
func SetOrientation(o Orientation) error {
	want := orientationUnspecified
	switch o {
	case Portrait:
		want = orientationPortrait
	case Landscape:
		want = orientationLandscape
	}
	return requestOrientation(want)
}

// LockOrientation freezes the screen whichever way up it is now.
func LockOrientation(on bool) error {
	if on {
		return requestOrientation(orientationLocked)
	}
	return requestOrientation(orientationUnspecified)
}

func requestOrientation(want int) error {
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			act := jni.Activity()
			cls, err := e.ClassOf(act)
			if err != nil {
				return err
			}
			defer e.Delete(cls.Object())
			m, err := e.Method(cls, "setRequestedOrientation", jni.Sig(jni.TVoid, jni.TInt))
			if err != nil {
				return err
			}
			return e.CallVoid(act, m, jni.Int(int32(want)))
		})
	}); e != nil {
		return e
	}
	return err
}

// SetKeepScreenOn stops the display dimming and locking while the app is in
// front. It is a window flag rather than a Java call, so it needs neither
// the UI thread nor the bridge.
func SetKeepScreenOn(on bool) {
	if on {
		SetWindowFlags(FlagKeepScreenOn, 0)
	} else {
		SetWindowFlags(0, FlagKeepScreenOn)
	}
}

// SetBrightness overrides the screen's brightness while this app is in
// front: 0 for as dark as the device allows, 1 for full, and a negative
// number to hand the decision back to the system.
func SetBrightness(level float32) error {
	if level > 1 {
		level = 1
	}
	if level < 0 {
		level = -1 // the platform's "use the user's setting"
	}
	var err error
	if e := RunOnUISync(func() {
		err = jni.Do(func(e *jni.Env) error {
			win, err := activityWindow(e)
			if err != nil {
				return err
			}
			winCls, err := e.Class("android/view/Window")
			if err != nil {
				return err
			}
			paramsType := jni.TClass("android/view/WindowManager$LayoutParams")
			get, err := e.Method(winCls, "getAttributes", jni.Sig(paramsType))
			if err != nil {
				return err
			}
			params, err := e.CallObject(win, get)
			if err != nil {
				return err
			}
			pCls, err := e.Class("android/view/WindowManager$LayoutParams")
			if err != nil {
				return err
			}
			f, err := e.Field(pCls, "screenBrightness", jni.TFloat)
			if err != nil {
				return err
			}
			if err := e.SetFloat(params, f, level); err != nil {
				return err
			}
			set, err := e.Method(winCls, "setAttributes", jni.Sig(jni.TVoid, paramsType))
			if err != nil {
				return err
			}
			return e.CallVoid(win, set, jni.Ref(params))
		})
	}); e != nil {
		return e
	}
	return err
}

// FontScale is how much bigger the user has asked text to be: 1 is the
// default, and it goes past 2 on a device set up for poor eyesight. An app
// that draws its own text is the only thing that can honour it.
func FontScale() (float32, error) {
	var scale float32 = 1
	err := jni.Do(func(e *jni.Env) error {
		act := jni.Activity()
		actCls, err := e.ClassOf(act)
		if err != nil {
			return err
		}
		defer e.Delete(actCls.Object())
		getRes, err := e.Method(actCls, "getResources",
			jni.Sig(jni.TClass("android/content/res/Resources")))
		if err != nil {
			return err
		}
		res, err := e.CallObject(act, getRes)
		if err != nil {
			return err
		}
		resCls, err := e.Class("android/content/res/Resources")
		if err != nil {
			return err
		}
		getCfg, err := e.Method(resCls, "getConfiguration",
			jni.Sig(jni.TClass("android/content/res/Configuration")))
		if err != nil {
			return err
		}
		cfg, err := e.CallObject(res, getCfg)
		if err != nil {
			return err
		}
		cfgCls, err := e.Class("android/content/res/Configuration")
		if err != nil {
			return err
		}
		f, err := e.Field(cfgCls, "fontScale", jni.TFloat)
		if err != nil {
			return err
		}
		scale, err = e.GetFloat(cfg, f)
		return err
	})
	return scale, err
}

// activityWindow is the activity's android.view.Window.
func activityWindow(e *jni.Env) (jni.Object, error) {
	act := jni.Activity()
	cls, err := e.ClassOf(act)
	if err != nil {
		return jni.Object{}, err
	}
	defer e.Delete(cls.Object())
	m, err := e.Method(cls, "getWindow", jni.Sig(jni.TClass("android/view/Window")))
	if err != nil {
		return jni.Object{}, err
	}
	return e.CallObject(act, m)
}

// decorView is the view the whole window is drawn into, which is what the
// system bars are asked about and told about.
func decorView(e *jni.Env) (jni.Object, error) {
	win, err := activityWindow(e)
	if err != nil {
		return jni.Object{}, err
	}
	winCls, err := e.Class("android/view/Window")
	if err != nil {
		return jni.Object{}, err
	}
	m, err := e.Method(winCls, "getDecorView", jni.Sig(jni.TClass("android/view/View")))
	if err != nil {
		return jni.Object{}, err
	}
	return e.CallObject(win, m)
}

// Context is the app's android.content.Context — which is the activity, and
// is the starting point for almost everything Android offers. It is a global
// reference the library owns and callers must not free.
func Context() jni.Object { return jni.Activity() }

// The names of the system services this library reaches for. Android has
// dozens more; these are the ones with a package above them.
const (
	ServiceClipboard    = "clipboard"
	ServiceNotification = "notification"
	ServiceConnectivity = "connectivity"
	ServiceVibrator     = "vibrator"
	ServiceWindow       = "window"
	ServicePower        = "power"
	ServiceInputMethod  = "input_method"
	ServiceDisplay      = "display"
	ServiceSensor       = "sensor"
	ServiceLocation     = "location"
)

// Service is a system service by name — the object every one of Android's
// managers is got from.
func Service(e *jni.Env, name string) (jni.Object, error) {
	js, err := e.String(name)
	if err != nil {
		return jni.Object{}, err
	}
	svc, err := e.Invoke(Context(), "getSystemService",
		jni.Sig(jni.TObject, jni.TString), jni.Ref(js))
	if err != nil {
		return jni.Object{}, err
	}
	if svc.IsNil() {
		return jni.Object{}, fmt.Errorf("app: this device has no %q service", name)
	}
	return svc, nil
}
