package svg

// pictureArrived is a picture landing in a drawing that was already painted:
// the canvas the last paint made is dropped, so the next paint shows the photo
// instead of the hole. The epoch is bumped with it, and a paint that started
// before the picture arrived checks it before caching the canvas it made — one
// that cached the hole now would never be repainted over it.
func (img *Image) pictureArrived() {
	img.mu.Lock()
	img.drawn = nil
	img.epoch++
	img.mu.Unlock()
}
