package css

// add folds one scanned layer and its size into the fold. It reports false
// when a repeat or position group cannot be read.
func (f *backgroundFold) add(lay backgroundLayer, size BackSize, ctx Units) bool {
	if len(lay.repToks) > 0 {
		rep, ok := backRepeatPair(lay.repToks)
		if !ok {
			return false
		}
		f.repeats = append(f.repeats, rep)
	}
	if len(lay.posToks) > 0 {
		p, ok := parseBackPosTokens(lay.posToks, ctx)
		if !ok {
			return false
		}
		f.poss = append(f.poss, p)
	}
	f.sizes = append(f.sizes, size)
	if lay.attSet {
		f.attaches = append(f.attaches, lay.att)
	}
	if len(lay.boxes) > 0 {
		f.origins = append(f.origins, lay.boxes[0])
		if len(lay.boxes) > 1 {
			f.clips = append(f.clips, lay.boxes[1])
		}
	}
	if lay.img.URL != "" || lay.img.Grad != nil {
		f.images = append(f.images, lay.img)
	}
	if lay.colorSet {
		f.color, f.colorSet = lay.color, true
	}
	return true
}
