package svg

// picturePlacement is where a picture's own pixels go inside its box: how much
// each axis of it is scaled by and which corner of the box it starts from. The
// fit decides the scale — the biggest that still fits for a meet, the biggest
// that covers for a slice, and one scale per axis for none — and the alignment
// decides what is left over of the box, a slice leaving a negative amount and
// pulling the picture past the edge it is not pinned to.
func picturePlacement(p *Picture, pw, ph float64) (sx, sy, ox, oy float64) {
	bx, by, bw, bh := p.Box[0], p.Box[1], p.Box[2], p.Box[3]
	if p.Fit == PictureFitNone {
		return bw / pw, bh / ph, bx, by
	}
	s := bw / pw
	if p.Fit == PictureFitSlice {
		if b := bh / ph; b > s {
			s = b
		}
	} else if b := bh / ph; b < s {
		s = b
	}
	fx, fy := p.Align.offsets()
	return s, s, bx + (bw-pw*s)*fx, by + (bh-ph*s)*fy
}
