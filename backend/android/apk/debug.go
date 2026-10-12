package apk

// Debug reports whether this is the shared debug key rather than a real one.
func (k Keystore) Debug() bool { return k.Alias == debugAlias && k.StorePass == debugPass }
