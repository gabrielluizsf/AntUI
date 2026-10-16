package svg

// fillPicture puts the picture an address points at into the slot the element
// was built with, and says whether the slot has a picture in it — now, for a
// file or a `data:` address, which answer in the time the read takes, or in a
// later frame, for a server, whose reading runs off this goroutine so a slow
// answer cannot hold up the frame. The slot is what the node paints from, and
// an element whose picture never arrives keeps the hole.
func (img *Image) fillPicture(target string, r pictureRequest) bool {
	if pic := img.pictures[target]; pic != nil {
		r.apply(img, pic)
		return true
	}
	if !fetchesOffFrame(target) {
		data, err := fetchPicture(target)
		return img.decodeInto(target, pictureRead{data: data, err: err}, r)
	}
	e := pictureBytesFor(target)
	e.mu.Lock()
	if e.done {
		read := e.read
		e.mu.Unlock()
		return img.decodeInto(target, read, r)
	}
	e.waiting = append(e.waiting, r.watcher(img))
	e.mu.Unlock()
	return true
}
