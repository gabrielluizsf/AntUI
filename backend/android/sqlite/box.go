//go:build android

package sqlite

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func box(e *jni.Env, v any) (jni.Object, error) {
	switch t := v.(type) {
	case nil:
		return jni.Object{}, nil
	case string:
		return e.String(t)
	case []byte:
		return e.NewBytes(t)
	case bool:
		n := int64(0)
		if t {
			n = 1
		}
		return e.Static("java/lang/Long", "valueOf",
			jni.Sig(jni.TClass("java/lang/Long"), jni.TLong), jni.Long(n))
	case int:
		return box(e, int64(t))
	case int32:
		return box(e, int64(t))
	case int64:
		return e.Static("java/lang/Long", "valueOf",
			jni.Sig(jni.TClass("java/lang/Long"), jni.TLong), jni.Long(t))
	case float32:
		return box(e, float64(t))
	case float64:
		return e.Static("java/lang/Double", "valueOf",
			jni.Sig(jni.TClass("java/lang/Double"), jni.TDouble), jni.Double(t))
	}
	return jni.Object{}, fmt.Errorf("sqlite: cannot bind %T", v)
}
