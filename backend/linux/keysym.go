//go:build linux || freebsd || openbsd || netbsd || dragonfly

package linux

// keysym resolves a keycode to the symbol it means with these modifiers held.
func (x *Driver) keysym(keycode byte, state uint16) uint32 {
	if x.keysyms == nil || keycode < x.minKeycode || keycode > x.maxKeycode {
		return 0
	}
	perCode := int(x.keysymsPerCode)
	base := (int(keycode) - int(x.minKeycode)) * perCode
	if base < 0 || base >= len(x.keysyms) {
		return 0
	}

	index := 0
	if state&stateAltGr != 0 && perCode > 2 {
		index = 2
	}
	if state&stateShift != 0 {
		index++
	}
	index = min(index, perCode-1)

	sym := x.keysyms[base+index]
	// An empty slot in the shifted column means the key has only one symbol.
	if sym == 0 && index > 0 {
		sym = x.keysyms[base]
	}

	// Caps lock only affects letters, and it inverts shift rather than
	// forcing upper case — which is why holding shift with caps on types
	// lower case, as it does everywhere else.
	switch {
	case state&stateLock != 0 && state&stateShift == 0 && sym >= 'a' && sym <= 'z':
		sym -= 32
	case state&stateLock != 0 && state&stateShift != 0 && sym >= 'A' && sym <= 'Z':
		sym += 32
	}
	return sym
}
