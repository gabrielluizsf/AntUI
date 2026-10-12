//go:build android

package speech

import (
	"github.com/gabrielluizsf/antui/backend/android/jni"
)

// readResults pulls the guesses out of the bundle the engine answers with.
func readResults(e *jni.Env, bundle jni.Object, partial bool) (Heard, bool) {
	if bundle.IsNil() {
		return Heard{}, false
	}
	key, err := e.String("results_recognition")
	if err != nil {
		return Heard{}, false
	}
	list, err := e.Invoke(bundle, "getStringArrayList",
		jni.Sig(jni.TClass("java/util/ArrayList"), jni.TString), jni.Ref(key))
	if err != nil || list.IsNil() {
		return Heard{}, false
	}
	n, err := e.InvokeInt(list, "size", jni.Sig(jni.TInt))
	if err != nil || n == 0 {
		return Heard{}, false
	}
	out := Heard{Partial: partial, Alternatives: make([]string, 0, n)}
	err = e.Frame(int(n)*4+16, func() error {
		for i := range int(n) {
			item, err := e.Invoke(list, "get", jni.Sig(jni.TObject, jni.TInt),
				jni.Int(int32(i)))
			if err != nil || item.IsNil() {
				continue
			}
			s, err := e.GoString(item)
			if err == nil && s != "" {
				out.Alternatives = append(out.Alternatives, s)
			}
		}
		return nil
	})
	if err != nil || len(out.Alternatives) == 0 {
		return Heard{}, false
	}
	out.Text = out.Alternatives[0]

	// The confidence scores, which many engines simply do not send.
	if ckey, err := e.String("confidence_scores"); err == nil {
		if scores, err := e.Invoke(bundle, "getFloatArray",
			jni.Sig(jni.TArray(jni.TFloat), jni.TString), jni.Ref(ckey)); err == nil &&
			!scores.IsNil() {
			if got, err := goFloats(e, scores); err == nil && len(got) > 0 {
				out.Confidence = got[0]
			}
		}
	}
	return out, true
}
