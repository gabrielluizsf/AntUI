package svg

// pictureRead is what fetching an address answered: the bytes, or the reason
// there are none. It is the one shape a fetch and everything waiting on it pass
// between them, whether the reading happened inside the drawing's own frame or
// off it.
type pictureRead struct {
	data []byte
	err  error
}

// decoded is the read turned into the frames the element draws, which is the
// only thing drawing an `<image>` needs from it. A failure is answered the same
// way whatever it was: as the reason the picture is not there, named by the
// error it carries.
func (r pictureRead) decoded() (*pictureFrames, error) {
	if r.err != nil {
		return nil, r.err
	}
	return decodePicture(r.data)
}
