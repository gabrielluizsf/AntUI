//go:build android

package jni

// Make builds an object: find the class, find the constructor, call it.
func (e *Env) Make(class, sig string, args ...Value) (Object, error) {
	c, err := e.Class(class)
	if err != nil {
		return Object{}, err
	}
	ctor, err := e.cachedMethod(c, "<init>", sig, kindMethod)
	if err != nil {
		return Object{}, err
	}
	return e.New(c, ctor, args...)
}
