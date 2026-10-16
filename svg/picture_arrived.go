package svg

// pictureArrived is a picture landing in a drawing that was already painted:
// the canvas the last paint made is dropped, so the next paint shows the photo
// instead of the hole. The epoch is bumped with it, and a paint that started
// before the picture arrived checks it before caching the canvas it made — one
// that cached the hole now would never be repainted over it. A picture that
// moves marks the drawing so its frames are chosen by the clock from here on.
func (img *Image) pictureArrived(moves bool) {
	img.mu.Lock()
	img.drawn = nil
	img.epoch++
	if moves {
		img.animated = true
	}
	img.mu.Unlock()
}
