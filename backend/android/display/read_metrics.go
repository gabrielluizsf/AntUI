//go:build android

package display

import "github.com/gabrielluizsf/antui/backend/android/jni"

// readMetrics reads the real size and density out of the display metrics.
func readMetrics(e *jni.Env, d jni.Object, out *Display) error {
	// getRealMetrics rather than getMetrics: the second reports the area an
	// app may use, which is the screen minus the system bars, and the first
	// reports the screen. This is the screen.
	metrics, err := e.Make("android/util/DisplayMetrics", jni.Sig(jni.TVoid))
	if err != nil {
		return err
	}
	if err := e.InvokeVoid(d, "getRealMetrics",
		jni.Sig(jni.TVoid, jni.TClass("android/util/DisplayMetrics")), jni.Ref(metrics)); err != nil {
		return err
	}
	cls, err := e.Class("android/util/DisplayMetrics")
	if err != nil {
		return err
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
			return err
		}
		v, err := e.GetInt(metrics, id)
		if err != nil {
			return err
		}
		*f.into = int(v)
	}
	return nil
}
