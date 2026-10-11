package template

// needsLayer reports whether the scratch layer has to be rebuilt: a window
// that grew or shrank since the last frame would otherwise draw into a stale
// surface.
func (cs *cssStyle) needsLayer() bool {
	w, h := cs.win.Width(), cs.win.Height()
	return cs.layer == nil || cs.layerW != w || cs.layerH != h
}
