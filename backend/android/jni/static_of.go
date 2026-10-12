//go:build android

package jni

func (e *Env) staticOf(class, name, sig string) (Class, Method, error) {
	c, err := e.Class(class)
	if err != nil {
		return Class{}, Method{}, err
	}
	m, err := e.cachedMethod(c, name, sig, kindStaticMethod)
	return c, m, err
}
