package css

// Font is the face a font-family list asks for at a weight and a slant, or
// nil when no @font-face in the sheet holds any of its names — the caller
// keeps drawing in the face it already had.
//
// The list is read name by name: the first name the sheet holds a face for
// draws, and every name after it is chained behind that face so the runes it
// lacks reach the next one that has them, ending at the face the program
// draws with by default. Inside one name CSS chooses: a face already leaning
// is looked for first when the style leans, and a straight one first when it
// does not; within that, the wanted weight and then the ones the spec tries
// from it.
//
// The answer is kept, so a frame asks once however many lines it draws — and
// so the chain it hands back ends at the default face the program had the
// first time the question was asked.
func (sh *Sheet) Font(family string, weight uint16, slanted bool) *FontFace {
	if sh == nil || family == "" {
		return nil
	}
	key := fontKey{family: family, weight: normWeight(weight), slanted: slanted}
	sh.fontMu.RLock()
	ff, seen := sh.fontCache[key]
	sh.fontMu.RUnlock()
	if seen {
		return ff
	}
	ff = sh.fontChain(family, weight, slanted)
	sh.fontMu.Lock()
	if sh.fontCache == nil {
		sh.fontCache = map[fontKey]*FontFace{}
	}
	sh.fontCache[key] = ff
	sh.fontMu.Unlock()
	return ff
}
