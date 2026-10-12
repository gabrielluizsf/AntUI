//go:build android

package jni

// InvokeString is [Env.Invoke] for a method that returns a String,
// read out as Go text.
func (e *Env) InvokeString(o Object, name, sig string, args ...Value) (string, error) {
	r, err := e.Invoke(o, name, sig, args...)
	if err != nil {
		return "", err
	}
	return e.GoString(r)
}
