//go:build android

package jni

// cachedMethod is Method or StaticMethod with the lookup remembered.
//
// A method id stays valid for as long as its class is loaded, and every
// class this package hands out is a global reference, so they never go. The
// cache is therefore permanent and that is correct rather than a leak: there
// are only so many methods in the platform.
func (e *Env) cachedMethod(c Class, name, sig string, kind uint8) (Method, error) {
	key := memberKey{uintptr(c.ref), name, sig, kind}
	memberMu.RLock()
	m, ok := methods[key]
	memberMu.RUnlock()
	if ok {
		return m, nil
	}
	var err error
	if kind == kindStaticMethod {
		m, err = e.StaticMethod(c, name, sig)
	} else {
		m, err = e.Method(c, name, sig)
	}
	if err != nil {
		return Method{}, err
	}
	memberMu.Lock()
	methods[key] = m
	memberMu.Unlock()
	return m, nil
}
