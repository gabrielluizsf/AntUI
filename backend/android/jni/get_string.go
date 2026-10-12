//go:build android

package jni

// GetString reads a String field and hands it back as Go text.
func (e *Env) GetString(o Object, f Field) (string, error) {
	r, err := e.GetObject(o, f)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
