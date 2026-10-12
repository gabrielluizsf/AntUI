//go:build android

package jni

func (e *Env) cachedField(c Class, name, sig string, kind uint8) (Field, error) {
	key := memberKey{uintptr(c.ref), name, sig, kind}
	memberMu.RLock()
	f, ok := fields[key]
	memberMu.RUnlock()
	if ok {
		return f, nil
	}
	var err error
	if kind == kindStaticField {
		f, err = e.StaticField(c, name, sig)
	} else {
		f, err = e.Field(c, name, sig)
	}
	if err != nil {
		return Field{}, err
	}
	memberMu.Lock()
	fields[key] = f
	memberMu.Unlock()
	return f, nil
}
