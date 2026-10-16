package svg

// Animates says the drawing holds a picture that moves with the clock — an
// animated GIF among its `<image>`s, whose frames are chosen as the drawing is
// painted. A canvas remembered from a paint a moment ago is a frame that has
// ended, so a caller that caches pictures asks this before keeping one.
func (img *Image) Animates() bool {
	if img == nil {
		return false
	}
	img.mu.Lock()
	defer img.mu.Unlock()
	return img.animated
}
