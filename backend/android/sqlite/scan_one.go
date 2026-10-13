//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

func (r *Rows) scanOne(e *jni.Env, i int32, dest any) error {
	kind, err := e.InvokeInt(r.cursor, "getType", jni.Sig(jni.TInt, jni.TInt), jni.Int(i))
	if err != nil {
		return err
	}
	c := column{e: e, cursor: r.cursor, i: i}
	return c.scan(kind, dest)
}
