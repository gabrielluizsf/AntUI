package antui

// pull fills a block of frames, and silences it when there is nothing to
// pull from. A backend hands its own buffer in.
func (a *Audio) pull(block []float32) {
	if a == nil || a.fill == nil {
		for i := range block {
			block[i] = 0
		}
		return
	}
	a.fill(block)
}
