//go:build android

package jni

// StaticIntOf is [Env.Static] for a class method that returns an int.
func (e *Env) StaticIntOf(class, name, sig string, args ...Value) (int32, error) {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return 0, err
	}
	return e.CallStaticInt(c, m, args...)
}
