package svg

// clearPatternPictures drops tiles made from a node tree that may have been
// restyled. The next paint rebuilds them from the refreshed content.
func (img *Image) clearPatternPictures() {
	for _, pattern := range img.patterns {
		pattern.mu.Lock()
		pattern.tiles = nil
		pattern.mu.Unlock()
	}
}
