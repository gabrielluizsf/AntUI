package css

// Keyframes returns the @keyframes block with the given name, or nil when
// the sheet never defined it.
func (sh *Sheet) Keyframes(name string) *Keyframes {
	if sh.keyframes == nil {
		return nil
	}
	return sh.keyframes[name]
}
