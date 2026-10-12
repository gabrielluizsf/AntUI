//go:build android

package jni

// Static calls a class method and hands back what it returned, looking the
// class and the method up once and keeping them — see [Env.Invoke], which is
// the same thing for a method on an object.
func (e *Env) Static(class, name, sig string, args ...Value) (Object, error) {
	c, m, err := e.staticOf(class, name, sig)
	if err != nil {
		return Object{}, err
	}
	return e.CallStaticObject(c, m, args...)
}
