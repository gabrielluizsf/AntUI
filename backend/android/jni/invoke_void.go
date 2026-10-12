//go:build android

package jni

// InvokeVoid, InvokeBool and the rest are Invoke for the other return types.
func (e *Env) InvokeVoid(o Object, name, sig string, args ...Value) error {
	m, err := e.method(o, name, sig)
	if err != nil {
		return err
	}
	return e.CallVoid(o, m, args...)
}
