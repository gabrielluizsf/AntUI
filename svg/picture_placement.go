package svg

// pictureFit is where a picture goes: the box it is drawn into, the fit it is
// given and the alignment that puts it in the box — everything
// [picturePlacement] asks of a slot, as one value so that a slot whose picture
// is still arriving off the frame can be read whole in the one moment it is
// locked.
type pictureFit struct {
	box   [4]float64
	fit   PictureFit
	align PictureAlign
}

// picturePlacement is where a picture's own pixels go inside its box: how much
// each axis of it is scaled by and which corner of the box it starts from. The
// fit decides the scale — the biggest that still fits for a meet, the biggest
// that covers for a slice, and one scale per axis for none — and the alignment
// decides what is left over of the box, a slice leaving a negative amount and
// pulling the picture past the edge it is not pinned to.
func picturePlacement(p pictureFit, pw, ph float64) (sx, sy, ox, oy float64) {
	bx, by, bw, bh := p.box[0], p.box[1], p.box[2], p.box[3]
	if p.fit == PictureFitNone {
		return bw / pw, bh / ph, bx, by
	}
	s := bw / pw
	if p.fit == PictureFitSlice {
		if b := bh / ph; b > s {
			s = b
		}
	} else if b := bh / ph; b < s {
		s = b
	}
	fx, fy := p.align.offsets()
	return s, s, bx + (bw-pw*s)*fx, by + (bh-ph*s)*fy
}
