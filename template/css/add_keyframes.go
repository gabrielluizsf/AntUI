package css

// addKeyframes stores one parsed block, a later definition of the same name
// replacing the earlier one the way CSS reads the last @keyframes it met.
func (sh *Sheet) addKeyframes(kf *Keyframes) {
	if sh.keyframes == nil {
		sh.keyframes = map[string]*Keyframes{}
	}
	sh.keyframes[kf.Name] = kf
}
