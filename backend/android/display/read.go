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
	if err := readMetrics(e, d, &out); err != nil {
		return out, err
	}

	// isHdr arrived in API 24. Below that the answer is no rather than an
	// error, because there were no HDR phones to say yes.
	if hdr, err := e.InvokeBool(d, "isHdr", jni.Sig(jni.TBool)); err == nil {
		out.HDR = hdr
	}
	return out, nil
}
