//go:build android

package speech

import "github.com/gabrielluizsf/antui/backend/android/jni"

// goFloats reads a Java float[].
//
// It is here rather than in the jni package because it is the only place in
// this library that has wanted one — the confidence scores a recogniser may
// or may not send.
func goFloats(e *jni.Env, arr jni.Object) ([]float32, error) {
	n, err := e.Len(arr)
	if err != nil || n == 0 {
		return nil, err
	}
	// Through java.lang.reflect.Array rather than a primitive-array call,
	// which the jni package does not expose for floats: n is one or two
	// here, so the cost of the reflection is nothing.
	out := make([]float32, 0, n)
	err = e.Frame(n+8, func() error {
		for i := range n {
			boxed, err := e.Static("java/lang/reflect/Array", "get",
				jni.Sig(jni.TObject, jni.TObject, jni.TInt), jni.Ref(arr), jni.Int(int32(i)))
			if err != nil || boxed.IsNil() {
				continue
			}
			v, err := e.InvokeFloat(boxed, "floatValue", jni.Sig(jni.TFloat))
			if err != nil {
				continue
			}
			out = append(out, v)
		}
		return nil
	})
	return out, err
}
