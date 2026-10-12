//go:build android

package jni

// StaticVoid is [Env.Static] for a class method that returns nothing.
func (e *Env) StaticVoid(class, name, sig string, args ...Value) error {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return err
	}
	return e.CallStaticVoid(c, m, args...)
}
