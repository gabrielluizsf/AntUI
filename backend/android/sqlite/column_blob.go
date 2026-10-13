//go:build android

package sqlite

import "github.com/gabrielluizsf/antui/backend/android/jni"

func (c column) blob() ([]byte, error) {
	o, err := c.e.Invoke(c.cursor, "getBlob",
		jni.Sig(jni.TArray(jni.TByte), jni.TInt), jni.Int(c.i))
	if err != nil || o.IsNil() {
		return nil, err
	}
	return c.e.GoBytes(o)
}
