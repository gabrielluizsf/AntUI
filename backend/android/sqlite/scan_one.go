//go:build android

package sqlite

import (
	"fmt"

	"github.com/gabrielluizsf/antui/backend/android/jni"
)

func (r *Rows) scanOne(e *jni.Env, i int32, dest any) error {
	kind, err := e.InvokeInt(r.cursor, "getType", jni.Sig(jni.TInt, jni.TInt), jni.Int(i))
	if err != nil {
		return err
	}
	text := func() (string, error) {
		return e.InvokeString(r.cursor, "getString", jni.Sig(jni.TString, jni.TInt), jni.Int(i))
	}
	whole := func() (int64, error) {
		return e.InvokeLong(r.cursor, "getLong", jni.Sig(jni.TLong, jni.TInt), jni.Int(i))
	}
	real := func() (float64, error) {
		v, err := e.InvokeFloat(r.cursor, "getFloat", jni.Sig(jni.TFloat, jni.TInt), jni.Int(i))
		return float64(v), err
	}
	blob := func() ([]byte, error) {
		o, err := e.Invoke(r.cursor, "getBlob",
			jni.Sig(jni.TArray(jni.TByte), jni.TInt), jni.Int(i))
		if err != nil || o.IsNil() {
			return nil, err
		}
		return e.GoBytes(o)
	}

	switch d := dest.(type) {
	case *string:
		if kind == typeNull {
			*d = ""
			return nil
		}
		v, err := text()
		*d = v
		return err
	case *int64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := whole()
		*d = v
		return err
	case *int:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := whole()
		*d = int(v)
		return err
	case *float64:
		if kind == typeNull {
			*d = 0
			return nil
		}
		v, err := real()
		*d = v
		return err
	case *bool:
		if kind == typeNull {
			*d = false
			return nil
		}
		v, err := whole()
		*d = v != 0
		return err
	case *[]byte:
		if kind == typeNull {
			*d = nil
			return nil
		}
		v, err := blob()
		*d = v
		return err
	case *any:
		switch kind {
		case typeNull:
			*d = nil
		case typeInt:
			v, err := whole()
			*d = v
			return err
		case typeReal:
			v, err := real()
			*d = v
			return err
		case typeBlob:
			v, err := blob()
			*d = v
			return err
		default:
			v, err := text()
			*d = v
			return err
		}
		return nil
	}
	return fmt.Errorf("sqlite: cannot scan into %T", dest)
}
