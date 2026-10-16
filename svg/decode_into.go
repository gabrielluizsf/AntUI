package svg

// decodeInto turns what a fetch answered into the picture the element draws,
// remembering it under the address for the other elements of the same drawing
// that name the same one. What could not be decoded is said and leaves the slot
// empty, the same as what could not be fetched.
func (img *Image) decodeInto(target string, read pictureRead, r pictureRequest) bool {
	pic, err := read.decoded(r.chain, target)
	if err != nil {
		r.unmet(err)
		return false
	}
	if img.pictures == nil {
		img.pictures = map[string]*pictureFrames{}
	}
	img.pictures[target] = pic
	r.apply(img, pic)
	return true
}
