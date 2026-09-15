//go:build android

package display

import (
	"github.com/gabrielluizsf/antui/backend/android/app"
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// Display is one screen: the built-in one, or a television the device is
// casting to, or a monitor it is plugged into.
type Display struct {
	// ID is 0 for the built-in screen and something else for the rest.
	ID   int
	Name string
	// Width and Height are the whole screen in pixels, including whatever
	// the system bars are covering.
	Width, Height int
	// Density is dots per inch as Android buckets it; 160 is one pixel per
	// device-independent pixel.
	Density int
	// Refresh is how many frames a second it shows. It is a float because a
	// panel that says 60 is usually running at 59.94 or 60.0000047, and
	// rounding it before the caller sees it loses the only interesting part.
	Refresh float32
	// HDR is whether it can show more than eight bits a channel. It is
	// false below API 24, where there was no way to ask.
	HDR bool
}

// Default is the screen the app is on.
func Default() (Display, error) {
	var out Display
	err := jni.Do(func(e *jni.Env) error {
		wm, err := app.Service(e, app.ServiceWindow)
		if err != nil {
			return err
		}
		d, err := e.Invoke(wm, "getDefaultDisplay",
			jni.Sig(jni.TClass("android/view/Display")))
		if err != nil {
			return err
		}
		if d.IsNil() {
			return errNoDisplay
		}
		out, err = read(e, d)
		return err
	})
	return out, err
}

// All is every screen the device can see, which is more than one while it is
// casting or plugged into a monitor.
func All() ([]Display, error) {
	var out []Display
	err := jni.Do(func(e *jni.Env) error {
		dm, err := app.Service(e, app.ServiceDisplay)
		if err != nil {
			return err
		}
		arr, err := e.Invoke(dm, "getDisplays",
			jni.Sig(jni.TArray(jni.TClass("android/view/Display"))))
		if err != nil || arr.IsNil() {
			return err
		}
		n, err := e.Len(arr)
		if err != nil {
			return err
		}
		out = make([]Display, 0, n)
		return e.Frame(n+8, func() error {
			for i := range n {
				d, err := e.Index(arr, i)
				if err != nil {
					return err
				}
				one, err := read(e, d)
				if err != nil {
					return err
				}
				out = append(out, one)
			}
			return nil
		})
	})
	return out, err
}

// Rotation is how far the screen is turned from the way the device was
// built: 0, 1, 2 or 3 for none, a quarter, a half and three quarters
// clockwise. It is what [antui/backend/android/sensor.Reading.ForDisplay] takes,
// and the reason it has to be asked for at all is that the sensors do not
// turn with the screen.
func Rotation() (int, error) {
	var out int
	err := jni.Do(func(e *jni.Env) error {
		wm, err := app.Service(e, app.ServiceWindow)
		if err != nil {
			return err
		}
		d, err := e.Invoke(wm, "getDefaultDisplay",
			jni.Sig(jni.TClass("android/view/Display")))
		if err != nil || d.IsNil() {
			return err
		}
		r, err := e.InvokeInt(d, "getRotation", jni.Sig(jni.TInt))
		out = int(r)
		return err
	})
	return out, err
}

// Refresh is how many frames a second the app's screen shows. It is what
// antui asks for when it wants to know the frame budget.
func Refresh() (float32, error) {
	d, err := Default()
	return d.Refresh, err
}

var errNoDisplay = errString("display: this device reports no screen")

type errString string

func (e errString) Error() string { return string(e) }

// read fills in everything about one Display.
func read(e *jni.Env, d jni.Object) (Display, error) {
	var out Display
	id, err := e.InvokeInt(d, "getDisplayId", jni.Sig(jni.TInt))
	if err != nil {
		return out, err
	}
	out.ID = int(id)
	if out.Name, err = e.InvokeString(d, "getName", jni.Sig(jni.TString)); err != nil {
		return out, err
	}
	if out.Refresh, err = e.InvokeFloat(d, "getRefreshRate", jni.Sig(jni.TFloat)); err != nil {
		return out, err
	}

	// getRealMetrics rather than getMetrics: the second reports the area an
	// app may use, which is the screen minus the system bars, and the first
	// reports the screen. This is the screen.
	metrics, err := e.Make("android/util/DisplayMetrics", jni.Sig(jni.TVoid))
	if err != nil {
		return out, err
	}
	if err := e.InvokeVoid(d, "getRealMetrics",
		jni.Sig(jni.TVoid, jni.TClass("android/util/DisplayMetrics")), jni.Ref(metrics)); err != nil {
		return out, err
	}
	cls, err := e.Class("android/util/DisplayMetrics")
	if err != nil {
		return out, err
	}
	for _, f := range []struct {
		name string
		into *int
	}{
		{"widthPixels", &out.Width},
		{"heightPixels", &out.Height},
		{"densityDpi", &out.Density},
	} {
		id, err := e.Field(cls, f.name, jni.TInt)
		if err != nil {
			return out, err
		}
		v, err := e.GetInt(metrics, id)
		if err != nil {
			return out, err
		}
		*f.into = int(v)
	}

	// isHdr arrived in API 24. Below that the answer is no rather than an
	// error, because there were no HDR phones to say yes.
	if hdr, err := e.InvokeBool(d, "isHdr", jni.Sig(jni.TBool)); err == nil {
		out.HDR = hdr
	}
	return out, nil
}
