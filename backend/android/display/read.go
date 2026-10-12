//go:build android

package display

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

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
